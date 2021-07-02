package pool

import (
	"fmt"
)

type Resource interface {
	Close()
	IsClosed() bool
}

// Resource Factory
type Factory func() (Resource, error)

// 3类ResourcePool相关的消息, 对应Pool的3种操作:
// Acquire - 申请一个Resource
// Release - 释放一个Resource
// Close   - 关闭当前Resource pool中的所有Resource(会调用resource.Close方法)
type releaseMessage struct {
	r Resource
}

type acquireMessage struct {
	resourceChan chan Resource
	errorChan    chan error
}

type closeMessage struct {
}

// 一个资源池
//
// 1. 有下面3个构造参数:
// factory - 构造一个Resource的函数
// idleCapacity - resource pool中最大Idle的资源数(超过了则调用Close方法关闭), 如果-1则没有限制
// maxResources - 最大资源数, 如果为-1则没有限制
//
// 2. 主要提供了下面3个接口:
// Acquire - 申请一个resource
// Release - 释放一个resource
// Close   - 关闭当前resource pool中的所有resource
//
//
// 补充:
// 当前版本的实现, 一个如果acquire的资源达到上限, 新的acquire请求会被放到activeWaits队列中, 队列中请求的释放依赖
// 上层的release操作, 如果没有上层release操作, 这些请求会一直被挂起.
type ResourcePool struct {
	factory       Factory // Resource Factory
	idleResources *Ring   // Idle Resource资源池(acquire会从Idle资源池中取出, release用完放回去)
	idleCapacity  int     // 最大Idle的资源数(超过了则调用Close方法关闭), 如果-1则没有限制
	maxResources  int     // 最大资源数, 如果为-1则没有限制
	numResources  int     // 当前总的资源数(包括已经在使用的和Idle的)

	acquireChan chan *acquireMessage // acquire资源
	releaseChan chan *releaseMessage // release资源
	closeChan   chan *closeMessage   // close资源池(会close所有资源)

	activeWaits []*acquireMessage // pending requests
}

func NewResourcePool(factory Factory, idleCapacity, maxResources int) *ResourcePool {
	resourcePool := &ResourcePool{
		factory:       factory,
		idleResources: NewRing(),
		idleCapacity:  idleCapacity,
		maxResources:  maxResources,
		numResources:  0,

		acquireChan: make(chan *acquireMessage),
		releaseChan: make(chan *releaseMessage, 1),
		closeChan:   make(chan *closeMessage, 1),

		activeWaits: make([]*acquireMessage, 0, 0),
	}

	go resourcePool.mux()

	return resourcePool
}

func (resourcePool *ResourcePool) Acquire() (resource Resource, err error) {
	acquireMsg := &acquireMessage{
		resourceChan: make(chan Resource),
		errorChan:    make(chan error),
	}
	resourcePool.acquireChan <- acquireMsg

	select {
	case resource = <-acquireMsg.resourceChan:
	case err = <-acquireMsg.errorChan:
	}

	return
}

func (resourcePool *ResourcePool) Release(resource Resource) {
	releaseMsg := &releaseMessage{
		r: resource,
	}
	resourcePool.releaseChan <- releaseMsg
}

// 关闭所有Resource
func (resourcePool *ResourcePool) Close() {
	resourcePool.closeChan <- &closeMessage{}
}

func (resourcePool *ResourcePool) NumResources() int {
	return resourcePool.numResources
}

// Message Loop
func (resourcePool *ResourcePool) mux() {
loop:
	for {
		select {
		case acquireMsg := <-resourcePool.acquireChan:
			// acquire resource
			resourcePool.acquire(acquireMsg)
		case releaseMsg := <-resourcePool.releaseChan:
			// release resource
			if len(resourcePool.activeWaits) != 0 {
				// 有request在waiting
				if !releaseMsg.r.IsClosed() {
					resourcePool.activeWaits[0].resourceChan <- releaseMsg.r
				} else {
					// new一个新Resource
					r, err := resourcePool.factory()
					if err != nil {
						resourcePool.numResources--
						resourcePool.activeWaits[0].errorChan <- err
					} else {
						resourcePool.activeWaits[0].resourceChan <- r
					}
				}

				resourcePool.activeWaits = resourcePool.activeWaits[1:]
			} else {
				// 没有request在waiting
				resourcePool.release(releaseMsg.r)
			}

		case _ = <-resourcePool.closeChan:
			// close resources
			break loop
		}
	}
	for !resourcePool.idleResources.Empty() {
		resourcePool.idleResources.Dequeue().Close()
	}
	for _, aw := range resourcePool.activeWaits {
		aw.errorChan <- fmt.Errorf("Resource pool closed")
	}
}

func (resourcePool *ResourcePool) acquire(acq *acquireMessage) {
	for !resourcePool.idleResources.Empty() {
		r := resourcePool.idleResources.Dequeue()
		if !r.IsClosed() {
			acq.resourceChan <- r
			return
		}

		// discard closed resources
		resourcePool.numResources--
	}

	// 资源池为空

	// 放入pending队列等待
	if resourcePool.maxResources != -1 && resourcePool.numResources >= resourcePool.maxResources {
		resourcePool.activeWaits = append(resourcePool.activeWaits, acq)
		return
	}

	// new一个新Resource
	r, err := resourcePool.factory()
	if err != nil {
		acq.errorChan <- err
	} else {
		resourcePool.numResources++
		acq.resourceChan <- r
	}

	return
}

// 如果资源已经Closed(), 则不会重新放入到Pool中
func (resourcePool *ResourcePool) release(resource Resource) {
	if resource == nil || resource.IsClosed() {
		// don't put it back in the pool.
		resourcePool.numResources--
		return
	}

	if resourcePool.idleCapacity != -1 && resourcePool.idleResources.Size() == resourcePool.idleCapacity {
		resource.Close()
		resourcePool.numResources--
		return
	}

	resourcePool.idleResources.Enqueue(resource)
}

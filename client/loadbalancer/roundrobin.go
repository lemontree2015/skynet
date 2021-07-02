package loadbalancer

import (
	"container/list"
	"errors"
	"monitor_server/skynet"
	"sync"
)

var (
	NoServices = errors.New("No services")
)

// 对一组ServiceInfo做LoadBalancer
type LoadBalancer struct {
	serviceMap  map[string]*list.Element
	serviceList *list.List
	current     *list.Element
	lock        *sync.RWMutex
}

func NewLoadBalancer(services []*skynet.ServiceInfo) *LoadBalancer {
	loadBalancer := &LoadBalancer{
		serviceMap:  make(map[string]*list.Element),
		serviceList: list.New(),
		current:     nil,
		lock:        new(sync.RWMutex),
	}

	for _, si := range services {
		loadBalancer.AddService(si)
	}

	return loadBalancer
}

func (loadBalancer *LoadBalancer) AddService(si *skynet.ServiceInfo) {
	loadBalancer.lock.Lock()
	defer loadBalancer.lock.Unlock()

	// 已经存在
	if _, ok := loadBalancer.serviceMap[si.ServiceUUID()]; ok {
		return
	}

	// 不存在
	var e *list.Element
	e = loadBalancer.serviceList.PushBack(si)     // update service list
	loadBalancer.serviceMap[si.ServiceUUID()] = e // update service map
}

func (loadBalancer *LoadBalancer) RemoveService(si *skynet.ServiceInfo) {
	loadBalancer.lock.Lock()
	defer loadBalancer.lock.Unlock()

	// 删除
	loadBalancer.serviceList.Remove(loadBalancer.serviceMap[si.ServiceUUID()]) // update service list
	delete(loadBalancer.serviceMap, si.ServiceUUID())                          // update service map

	// current should be nil if we have no services
	if loadBalancer.serviceList.Len() == 0 {
		loadBalancer.current = nil
	}
}

// [20160219 Bug Fix]:
// 一个ServiceClient对应一个Loadbalancer, 当这个Loadbalancer在多个goroutine中被竞争的时候, 有可能同一个goroutine
// Choose出来的ServiceInfo和前一个是相同的(这在一个Loadbalancer中存在Deal ServiceInfo的时候会导致Retry失效), 增加一个
// services参数, 可以保证返回的结果一定不在services中.
//
// 参数:
// 如果si为nil, 直接执行正常逻辑
// 如果si不为nil, 保证返回的结果一定不在services中
func (loadBalancer *LoadBalancer) Choose(services []*skynet.ServiceInfo) (*skynet.ServiceInfo, error) {
	loadBalancer.lock.Lock()
	defer loadBalancer.lock.Unlock()

	// LoadBalancer没有数据
	if loadBalancer.serviceList.Len() == 0 {
		return nil, NoServices
	}

	// LoadBalancer有数据
	if services == nil || len(services) == 0 {
		// 不需要过滤
		return loadBalancer.next(), nil
	} else {
		// 需要过滤
		var si *skynet.ServiceInfo
		for i := 0; i < loadBalancer.serviceList.Len(); i++ {
			si = loadBalancer.next()
			canReturn := true
			for _, v := range services {
				if si.Equal(v) {
					canReturn = false // 在列表中则不能返回
				}
			}

			if canReturn {
				return si, nil
			}
		}

		return nil, NoServices
	}
}

// 前提是:
// loadBalancer.serviceList.Len() > 0
func (loadBalancer *LoadBalancer) next() *skynet.ServiceInfo {
	if loadBalancer.current == nil {
		loadBalancer.current = loadBalancer.serviceList.Front()
		return loadBalancer.current.Value.(*skynet.ServiceInfo)
	}

	loadBalancer.current = loadBalancer.current.Next() // 下一个list.Element
	if loadBalancer.current == nil {
		loadBalancer.current = loadBalancer.serviceList.Front()
	}

	return loadBalancer.current.Value.(*skynet.ServiceInfo)
}

package pool

import (
	"github.com/golang/glog"
	"github.com/lemontree2015/skynet"
	"github.com/lemontree2015/skynet/client/conn"
	"github.com/lemontree2015/skynet/config"
	"github.com/lemontree2015/skynet/cron"
	"github.com/lemontree2015/skynet/misc"
	"runtime/debug"
	"sync"
)

// 全局的Connection Pools(ServicePool的集合)
var globalConnPool *ConnPool = NewConnPool()

func GetService(si *skynet.ServiceInfo) *ServicePool {
	return globalConnPool.GetService(si)
}

func RemoveService(si *skynet.ServiceInfo) {
	globalConnPool.RemoveService(si)
}

func Acquire(si *skynet.ServiceInfo) (*conn.Connection, error) {
	return globalConnPool.Acquire(si)
}

func Release(c *conn.Connection) {
	globalConnPool.Release(c)
}

func Close() {
	globalConnPool.Close()
}

func NumConnections() int {
	return globalConnPool.NumConnections()
}

func NumServices() int {
	return globalConnPool.NumServices()
}

type ConnPool struct {
	servicePools map[string]*ServicePool
	cronEvery    *cron.CronEvery
	lock         *sync.RWMutex
}

func NewConnPool() *ConnPool {
	connPool := &ConnPool{
		servicePools: make(map[string]*ServicePool),
		lock:         new(sync.RWMutex),
	}
	connPool.cronEvery = cron.NewCronEvery(config.PoolGCInterval(), connPool.cronEveryFun)
	return connPool
}

// 备注:
// 这个函数一定会成功
func (pool *ConnPool) GetService(si *skynet.ServiceInfo) *ServicePool {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if servicePool, ok := pool.servicePools[si.ServiceUUID()]; ok {
		return servicePool
	} else {
		servicePool := NewServicePool(si) // 这个函数一定会成功
		pool.servicePools[si.ServiceUUID()] = servicePool
		return servicePool
	}
}

func (pool *ConnPool) RemoveService(si *skynet.ServiceInfo) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if servicePool, ok := pool.servicePools[si.ServiceUUID()]; ok {
		servicePool.Close()
	}

	delete(pool.servicePools, si.ServiceUUID())
}

func (pool *ConnPool) Acquire(si *skynet.ServiceInfo) (*conn.Connection, error) {
	servicePool := pool.GetService(si)                // 一定成功
	servicePool.touchTimestamp = misc.UnixTimestamp() // Touch(防止被GC掉)

	if r, err := servicePool.pool.Acquire(); err != nil {
		return nil, err
	} else {
		return r.(*conn.Connection), nil
	}
}

func (pool *ConnPool) Release(c *conn.Connection) {
	if c == nil {
		return
	}

	servicePool := pool.GetService(c.ServiceInfo())
	servicePool.pool.Release(c)
}

func (pool *ConnPool) Close() {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	for k, sp := range pool.servicePools {
		sp.Close()
		delete(pool.servicePools, k)
	}

}

func (pool *ConnPool) cronEveryFun() {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	defer func() {
		if r := recover(); r != nil {
			glog.Infof("ServiceClient.cronEveryFun: %v, DEBUG.STACK=%v", r, string(debug.Stack()))
			return
		}
	}()

	timestamp := misc.UnixTimestamp()
	for k, sp := range pool.servicePools {
		if sp.service.Name == "Monitor Service" {
			// Monitor Service不参与GC
			glog.Infof("GC Pool(Normal): touchTimestamp=%v, delta=%v, service=%v",
				sp.TouchTimestamp(), timestamp-sp.TouchTimestamp(), sp.service.ServiceUUID())
		} else {
			if sp.CanGC() {
				glog.Infof("GC Pool(Delete): touchTimestamp=%v, delta=%v, service=%v",
					sp.TouchTimestamp(), timestamp-sp.TouchTimestamp(), sp.service.ServiceUUID())
				sp.Close()
				delete(pool.servicePools, k)
			} else {
				glog.Infof("GC Pool(Normal): touchTimestamp=%v, delta=%v, service=%v",
					sp.TouchTimestamp(), timestamp-sp.TouchTimestamp(), sp.service.ServiceUUID())
			}
		}
	}
}

func (pool *ConnPool) NumConnections() (count int) {
	pool.lock.RLock()
	defer pool.lock.RUnlock()

	for _, sp := range pool.servicePools {
		count += sp.NumResources()
	}

	return count
}

func (pool *ConnPool) NumServices() int {
	pool.lock.RLock()
	defer pool.lock.RUnlock()

	return len(pool.servicePools)
}

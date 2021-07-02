package pool

import (
	"monitor_server/skynet"
	"monitor_server/skynet/client/conn"
	"monitor_server/skynet/config"
	"monitor_server/skynet/misc"
)

// 每个ServiceInfo对应一个ServicePool

type ServicePool struct {
	service        *skynet.ServiceInfo
	pool           *ResourcePool
	touchTimestamp int64
}

// 备注:
// 这个函数一定会成功
func NewServicePool(si *skynet.ServiceInfo) *ServicePool {
	pool := NewResourcePool(
		func() (Resource, error) {
			return conn.NewConnection(si, "tcp", config.ClientRPCDialTimeout(si.Name, si.Version))
		}, config.ClientConnIdle(si.Name, si.Version), config.ClientConnMax(si.Name, si.Version))

	return &ServicePool{
		service:        si,
		pool:           pool,
		touchTimestamp: misc.UnixTimestamp(),
	}

}

func (servicePool *ServicePool) Close() {
	servicePool.pool.Close()
}

func (servicePool *ServicePool) NumResources() int {
	return servicePool.pool.NumResources()
}

func (servicePool *ServicePool) TouchTimestamp() int64 {
	return servicePool.touchTimestamp
}

func (servicePool *ServicePool) CanGC() bool {
	return misc.UnixTimestamp()-servicePool.touchTimestamp > int64(config.PoolGCTimeout())
}

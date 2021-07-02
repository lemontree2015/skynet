package client

import (
	"errors"
	"github.com/golang/glog"
	"monitor_server/skynet"
	"monitor_server/skynet/client/loadbalancer"
	"monitor_server/skynet/client/pool"
	"monitor_server/skynet/config"
	"monitor_server/skynet/cron"
	"runtime/debug"
	"sync"
)

var (
	ServiceClientClosed = errors.New("Service client shutdown")
	MaxRetryCount       = errors.New("Max Retry Count")
)

type ServiceClient struct {
	loadBalancer *loadbalancer.LoadBalancer // 目标ServiceInfo的LoadBalancer
	criteria     *skynet.Criteria           // 可以匹配目标ServiceInfo
	closed       bool                       // ServiceClient是否已经Close
	cronEvery    *cron.CronEvery
	lock         *sync.RWMutex
}

func NewServiceClient(criteria *skynet.Criteria) *ServiceClient {
	serviceClient := &ServiceClient{
		loadBalancer: nil,
		criteria:     criteria,
		closed:       true,
		lock:         new(sync.RWMutex),
	}
	serviceClient.cronEvery = cron.NewCronEvery(config.CientSyncInterval(), serviceClient.cronEveryFun)
	return serviceClient
}

func (sc *ServiceClient) UpdateServices(services []*skynet.ServiceInfo) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	newLoadbalancer := loadbalancer.NewLoadBalancer([]*skynet.ServiceInfo{})
	for _, si := range services {
		if sc.criteria.Matches(si) {
			newLoadbalancer.AddService(si)
		}
	}

	sc.loadBalancer = newLoadbalancer
	sc.closed = false // open close flag
}

func (sc *ServiceClient) cronEveryFun() {
	if sc.closed {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			glog.Infof("ServiceClient.cronEveryFun: %v, DEBUG.STACK=%v", r, string(debug.Stack()))
			return
		}
	}()

	// Monitor Client特殊处理
	//
	// New出来之后, 永远不需要再Sync
	if len(sc.criteria.ServiceCriterias) > 0 {
		if sc.criteria.ServiceCriterias[0].Name == "Monitor Service" {
			return
		}
	}

	// 从Monitor获取Service
	if services, runTime, err := GetServices(sc.criteria); err == nil {
		if int(runTime) > config.MonitorTrustTime() {
			// Monitor是可信任的
			sc.UpdateServices(services)
		}
	}
}

func (sc *ServiceClient) Close() {
	sc.closed = true
	sc.loadBalancer = nil
	sc.cronEvery.Stop()
}

func (sc *ServiceClient) IsClosed() bool {
	return sc.IsClosed()
}

func (sc *ServiceClient) Send(fn string, in interface{}, out interface{}) error {
	return sc.trySend(config.ClientRPCRetry(), fn, in, out, nil)
}

func (sc *ServiceClient) trySend(retry int, fn string, in interface{}, out interface{}, filterServices []*skynet.ServiceInfo) error {
	if retry <= 0 {
		return MaxRetryCount
	}

	if sc.closed {
		return ServiceClientClosed
	}

	if sc.loadBalancer == nil {
		return ServiceClientClosed
	}

	if si, err := sc.loadBalancer.Choose(filterServices); err == nil {
		conn, err := pool.Acquire(si)
		defer pool.Release(conn)
		if err == nil {
			if err = conn.Send(fn, in, out); err == nil {
				return nil
			} else {
				// TCP层面的发送失败(Retry)
				// 1. connect关闭
				// 2. timeout超时
				return sc.trySend(retry-1, fn, in, out, append(filterServices, si))
			}
		} else {
			// 获取Connection失败(Retry)
			return sc.trySend(retry-1, fn, in, out, append(filterServices, si))
		}
	} else {
		// 获取ServiceInfo失败(Retry)
		return sc.trySend(retry-1, fn, in, out, filterServices) // Bug Fix: 不能把nil append到filterServices
	}

	return ServiceClientClosed
}

package client

import (
	"github.com/lemontree2015/skynet/config"
	"github.com/lemontree2015/skynet/cron"
	"github.com/lemontree2015/skynet/logger"
	"runtime/debug"
	"sync"
)

var (
	globalClientManager *ClientManager
)

type ClientManager struct {
	clients   map[string]*ServiceClient // key = ServiceName-Version
	cronEvery *cron.CronEvery
	lock      *sync.RWMutex
}

func NewClientManager() *ClientManager {
	clientManager := &ClientManager{
		clients: make(map[string]*ServiceClient),
		lock:    new(sync.RWMutex),
	}
	clientManager.cronEvery = cron.NewCronEvery(config.CientSyncInterval(), clientManager.cronEveryFun)
	return clientManager
}

func (manager *ClientManager) GetClient(serviceName, version string) *ServiceClient {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if c, ok := manager.clients[serviceName+"-"+version]; ok {
		return c
	} else {
		c := NewClient(serviceName, version)
		manager.clients[serviceName+"-"+version] = c
		return c
	}
}

func (manager *ClientManager) DelClient(serviceName, version string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	delete(manager.clients, serviceName+"-"+version)
}

func (manager *ClientManager) DelClientByServiceKey(serviceKey string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if _, ok := manager.clients[serviceKey]; ok {
		delete(manager.clients, serviceKey)
	}
}

func (manager *ClientManager) cronEveryFun() {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Infof("ClientManager.cronEveryFun: %v, DEBUG.STACK=%v", r, string(debug.Stack()))
			return
		}
	}()

	for serviceKey, cli := range manager.clients {
		if cli.IsClosed() {
			manager.DelClientByServiceKey(serviceKey)
		}
	}
}

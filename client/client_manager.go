package client

import (
	"sync"
)

var (
	globalClientManager *ClientManager
)

type ClientManager struct {
	clients map[string]*ServiceClient // key = ServiceName-Version
	lock    *sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]*ServiceClient),
		lock:    new(sync.RWMutex),
	}
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

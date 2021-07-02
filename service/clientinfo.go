package service

import (
	"net"
	"sync"
)

// 保存客户端信息

type clientInfo struct {
	ClientUUID string
	Address    net.Addr // Client Addr
}

type clientInfoManager struct {
	clientInfos map[string]*clientInfo // key = clientInfo.ClientUUID
	lock        *sync.RWMutex
}

func NewclientInfoManager() *clientInfoManager {
	return &clientInfoManager{
		clientInfos: make(map[string]*clientInfo),
		lock:        new(sync.RWMutex),
	}
}

func (manager *clientInfoManager) Set(client *clientInfo) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	if client != nil {
		manager.clientInfos[client.ClientUUID] = client
	}
}

func (manager *clientInfoManager) Get(clientUUID string) *clientInfo {
	manager.lock.RLock()
	defer manager.lock.RUnlock()

	return manager.clientInfos[clientUUID]
}

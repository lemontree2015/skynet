package client

import (
	"monitor_server/skynet"
)

func init() {
	globalMonitorClient = NewMonitorClient()
	globalClientManager = NewClientManager()
}

// 会重用全局的ServiceClient
//
// 备注:
// 一定是成功的
func GetClient(name, version string) *ServiceClient {
	return globalClientManager.GetClient(name, version)
}

// 创建一个新的ServiceClient
//
// 备注:
// 一定是成功的
func NewClient(name, version string) *ServiceClient {
	// 创建criteria
	criteria := skynet.NewCriteria()
	criteria.AddService(name, version)

	// 从Monitor Service获取符合条件的Services
	//
	// 注意:
	// New Service的时候, 我们认为Monitor Service是可信的(不需要判断运行时间),
	// 这里和Sync的逻辑不同
	if services, _, err := GetServices(criteria); err == nil {
		serviceClient := NewServiceClient(criteria)
		serviceClient.UpdateServices(services)
		return serviceClient
	} else {
		serviceClient := NewServiceClient(criteria)
		serviceClient.UpdateServices([]*skynet.ServiceInfo{})
		return serviceClient
	}
}

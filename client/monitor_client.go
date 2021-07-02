package client

import (
	"fmt"
	"github.com/lemontree2015/skynet"
)

var globalMonitorClient *ServiceClient

func NewMonitorClient() *ServiceClient {
	criteria := skynet.NewCriteria()
	criteria.AddService("Monitor Service", "1.0.0")

	// Monitor Service只限定一个port, 所以我们可以这样new出来直接使用
	cli := NewServiceClient(criteria)
	monitorService := skynet.NewServiceInfo("Monitor Service", "1.0.0")
	cli.UpdateServices([]*skynet.ServiceInfo{monitorService})

	return cli
}

/////////////////
// RPC Register
/////////////////

type RPCMonitorRegisterIn struct {
	Si *skynet.ServiceInfo
}

type RPCMonitorRegisterOut struct {
	Code string
}

func RegisterServices(si *skynet.ServiceInfo) error {
	// 构造参数
	in := &RPCMonitorRegisterIn{
		Si: si,
	}
	out := &RPCMonitorRegisterOut{}

	// RPC请求
	if err := globalMonitorClient.Send("Register", in, out); err == nil && out != nil {
		if out.Code == "success" {
			return nil
		} else {
			return fmt.Errorf(out.Code)
		}
	} else {
		return err
	}
}

///////////////////
// RPC UnRegister
///////////////////

type RPCMonitorUnRegisterIn struct {
	Si *skynet.ServiceInfo
}

type RPCMonitorUnRegisterOut struct {
	Code string
}

func UnRegisterServices(si *skynet.ServiceInfo) error {
	// 构造参数
	in := &RPCMonitorUnRegisterIn{
		Si: si,
	}
	out := &RPCMonitorUnRegisterOut{}

	// RPC请求
	if err := globalMonitorClient.Send("UnRegister", in, out); err == nil && out != nil {
		if out.Code == "success" {
			return nil
		} else {
			return fmt.Errorf(out.Code)
		}
	} else {
		return err
	}
}

/////////////
// RPC Get
/////////////

type RPCMonitorGetIn struct {
	Criteria *skynet.Criteria
}

type RPCMonitorGetOut struct {
	Code     string
	Services []*skynet.ServiceInfo
	RunTime  int64 // 运行时间(Seconds)
}

func GetServices(criteria *skynet.Criteria) ([]*skynet.ServiceInfo, int64, error) {
	// 构造参数
	in := &RPCMonitorGetIn{
		Criteria: criteria,
	}
	out := &RPCMonitorGetOut{}

	// RPC请求
	if err := globalMonitorClient.Send("Get", in, out); err == nil && out != nil {
		if out.Code == "success" {
			return out.Services, out.RunTime, nil
		} else {
			return nil, 0, fmt.Errorf(out.Code)
		}
	} else {
		return nil, 0, err
	}
}

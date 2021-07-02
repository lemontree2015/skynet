package skynet

import (
	"fmt"
	"github.com/lemontree2015/skynet/config"
	"strings"
)

type ServiceKey struct {
	Name    string
	Version string
}

func NewServiceKey(name, version string) *ServiceKey {
	return &ServiceKey{
		Name:    name,
		Version: version,
	}
}

func (sk *ServiceKey) Equal(other *ServiceKey) bool {
	if sk.Name == other.Name && sk.Version == other.Version {
		return true
	}

	return false
}

func (sk *ServiceKey) String() string {
	return fmt.Sprintf("%v-%v", sk.Name, sk.Version)
}

// 一个Service的标识
type ServiceInfo struct {
	*ServiceKey
	InstanceUUID string // Instance UUID (标识这个ServiceInfo属于哪个Instance)
	//Name         string // 构造的时候指定
	//Version      string // 构造的时候指定
	Region string // region

	// 每个ServiceInfo都会有一个Listen的BindAddr
	Addr *BindAddr // 从配置文件中读取: host, service.port.min, service.port.max

	// 辅助属性
	LastRegisteredTimestamp int64 // 上次上报的时间戳
}

func NewServiceInfo(name, version string) *ServiceInfo {
	var host, region string
	var minPort, maxPort int

	// 注意:
	// 先读取Service Level的配置, 如果没有, 使用DEFAULT Level的配置
	if h, err := config.String(name, version, "host"); err == nil {
		host = h
	} else {
		if h, err := config.DefaultString("host"); err == nil {
			host = h
		} else {
			panic(fmt.Errorf("NewServiceInfo Error:%v, %v", name, version))
		}
	}

	if r, err := config.String(name, version, "region"); err == nil {
		region = r
	} else {
		if r, err := config.DefaultString("region"); err == nil {
			region = r
		} else {
			panic(fmt.Errorf("NewServiceInfo Error:%v, %v", name, version))
		}
	}

	if p, err := config.Int(name, version, "service.port.min"); err == nil {
		minPort = p
	} else {
		if p, err := config.DefaultInt("service.port.min"); err == nil {
			minPort = p
		} else {
			panic(fmt.Errorf("NewServiceInfo Error:%v, %v", name, version))
		}
	}

	if p, err := config.Int(name, version, "service.port.max"); err == nil {
		maxPort = p
	} else {
		if p, err := config.DefaultInt("service.port.max"); err == nil {
			maxPort = p
		} else {
			panic(fmt.Errorf("NewServiceInfo Error:%v, %v", name, version))
		}
	}

	addrStr := fmt.Sprintf("%v:%v-%v", host, minPort, maxPort)
	bindAddr, err := NewBindAddr(addrStr)

	if err != nil {
		panic(fmt.Errorf("NewServiceInfo Error:%v, %v", name, version))
	}

	return &ServiceInfo{
		ServiceKey: NewServiceKey(name, version),
		//Name:         name,
		//Version:      version,
		InstanceUUID:            config.InstanceUUID(),
		Region:                  region,
		Addr:                    bindAddr,
		LastRegisteredTimestamp: 0,
	}
}

func (si *ServiceInfo) Equal(otherSi *ServiceInfo) bool {
	if si.InstanceUUID == otherSi.InstanceUUID &&
		si.Name == otherSi.Name &&
		si.Version == otherSi.Version &&
		si.Region == otherSi.Region &&
		si.Addr.Equal(otherSi.Addr) {
		return true
	}

	return false
}

func (si *ServiceInfo) ServiceUUID() string {
	return fmt.Sprintf("%v-%v:%v:%v:%v", si.Name, si.Version, si.InstanceUUID, si.Region, si.Addr.String())
}

func (si *ServiceInfo) String() string {
	return fmt.Sprintf("%v-%v:%v:%v:%v", si.Name, si.Version, si.InstanceUUID, si.Region, si.Addr.String())
}

func ParseServiceNV(str string) (serviceName, version string, err error) {
	if index := strings.Index(str, "_"); index == -1 {
		return "", "", fmt.Errorf("ParseServiceNV Error: =%v", str)
	} else {
		return str[0:index], str[index+1:], nil
	}

	panic(fmt.Errorf("Logic Error"))
}

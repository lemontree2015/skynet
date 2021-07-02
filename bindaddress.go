package skynet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

var portLock *sync.RWMutex = new(sync.RWMutex) // 全局的Port Mutex, 防止端口申请的时候冲突

// net.Addr的增强版(可递增查找可用端口)
type BindAddr struct {
	IPAddress string // xxx.xxx.xxx.xxx
	Port      int
	MinPort   int
	MaxPort   int
}

// addr格式:
// 0.0.0.0:1000-2000  {Port = 1000, MinPort = 1000, MaxPort = 2000}
func NewBindAddr(addr string) (*BindAddr, error) {
	// addr为空
	if addr == "" {
		return nil, fmt.Errorf("NewBindAddr Failed, Invalid Addr: %v", addr)
	}

	// addr格式非法, 没有:
	split := strings.Index(addr, ":")
	if split == -1 {
		return nil, fmt.Errorf("NewBindAddr Failed, Invalid Addr: %v", addr)
	}

	ipAddress := addr[:split]
	if ipAddress == "" {
		return nil, fmt.Errorf("NewBindAddr Failed, Invalid Addr: %v", addr)
	}

	portStr := addr[split+1:]

	// addr格式非法, 没有-
	rindex := strings.Index(portStr, "-")
	if rindex == -1 {
		return nil, fmt.Errorf("NewBindAddr Failed, Invalid Addr: %v", addr)
	}

	minPortStr := portStr[:rindex]
	maxPortStr := portStr[rindex+1:]

	minPort, err := strconv.Atoi(minPortStr)
	if err != nil {
		return nil, fmt.Errorf("NewBindAddr Failed, Invalid Addr: %v", addr)
	}

	maxPort, err := strconv.Atoi(maxPortStr)
	if err != nil {
		return nil, fmt.Errorf("NewBindAddr Failed, Invalid Addr: %v", addr)
	}

	return &BindAddr{
		IPAddress: ipAddress,
		Port:      minPort,
		MinPort:   minPort,
		MaxPort:   maxPort,
	}, nil
}

func (addr *BindAddr) Listen() (*net.TCPListener, error) {
	// 会全局Lock, 防止port冲突
	portLock.Lock()
	defer portLock.Unlock()

	for {
		var tcpAddr *net.TCPAddr
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr.String())

		if err != nil {
			panic(err)
		}

		if listener, err := net.ListenTCP("tcp", tcpAddr); err == nil {
			// Listen成功
			return listener, nil
		}

		if addr.Port < addr.MaxPort {
			addr.Port++
		} else {
			return nil, fmt.Errorf("BindAddr.Listen Failed, Addr: %v", addr)
		}

	}

	return nil, fmt.Errorf("BindAddr.Listen Failed, Addr: %v", addr)
}

func (addr *BindAddr) Equal(otherAddr *BindAddr) bool {
	if addr.IPAddress == otherAddr.IPAddress &&
		addr.Port == otherAddr.Port &&
		addr.MinPort == otherAddr.MinPort &&
		addr.MaxPort == otherAddr.MaxPort {
		return true
	}

	return false
}

func (addr *BindAddr) String() string {
	return fmt.Sprintf("%v:%v", addr.IPAddress, addr.Port)
}

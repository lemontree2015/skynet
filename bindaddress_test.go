package skynet

import (
	"testing"
)

func TestBindAddrValid(t *testing.T) {
	// 合法性测试
	addr1, err := NewBindAddr("127.0.0.1:1000-2000")
	if err != nil ||
		addr1.IPAddress != "127.0.0.1" ||
		addr1.Port != 1000 ||
		addr1.MinPort != 1000 ||
		addr1.MaxPort != 2000 {
		t.Fatal("TestBindAddrValid Failed")
	}
}

func TestBindAddrInValid(t *testing.T) {
	// 非法性测试
	addr1, err := NewBindAddr("1000")
	if err == nil ||
		addr1 != nil {
		t.Fatal("TestBindAddrInValid Failed")
	}

	addr2, err := NewBindAddr("127.0.0.1:1000")
	if err == nil ||
		addr2 != nil {
		t.Fatal("TestBindAddrInValid Failed")
	}
}

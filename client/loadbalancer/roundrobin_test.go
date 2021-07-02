package loadbalancer

import (
	"github.com/lemontree2015/skynet"
	"testing"
)

func TestRoundRobinBasic(t *testing.T) {
	si1 := skynet.NewServiceInfo("name1", "1.0.0")
	si2 := skynet.NewServiceInfo("name2", "1.0.0")
	si3 := skynet.NewServiceInfo("name3", "1.0.0")
	si4 := skynet.NewServiceInfo("name4", "1.0.0")
	balancer := NewLoadBalancer([]*skynet.ServiceInfo{si1, si2, si3})
	balancer.AddService(si4)

	var si *skynet.ServiceInfo
	si, _ = balancer.Choose()
	if !si.Equal(si1) {
		t.Fatalf("TestRoundRobinBasic Failed: %v, %v", si, si1)
	}

	si, _ = balancer.Choose()
	if !si.Equal(si2) {
		t.Fatalf("TestRoundRobinBasic Failed: %v, %v", si, si2)
	}

	si, _ = balancer.Choose()
	if !si.Equal(si3) {
		t.Fatalf("TestRoundRobinBasic Failed: %v, %v", si, si3)
	}

	si, _ = balancer.Choose()
	if !si.Equal(si4) {
		t.Fatalf("TestRoundRobinBasic Failed: %v, %v", si, si4)
	}

	si, _ = balancer.Choose()
	if !si.Equal(si1) {
		t.Fatalf("TestRoundRobinBasic Failed: %v, %v", si, si1)
	}

	si, _ = balancer.Choose()
	if !si.Equal(si2) {
		t.Fatalf("TestRoundRobinBasic Failed: %v, %v", si, si2)
	}
}

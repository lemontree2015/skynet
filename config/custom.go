package config

import (
	"fmt"
	"time"
)

func MonitorTrustTime() int {
	if v, err := DefaultInt("monitor.trust.time"); err == nil {
		return v
	} else {
		panic(fmt.Errorf("monitor.trust.time Error"))
	}
}

func CientSyncInterval() time.Duration {
	if v, err := DefaultInt("client.sync.interval"); err == nil {
		return time.Duration(v) * time.Second
	} else {
		panic(fmt.Errorf("client.sync.interval Error"))
	}
}

func PoolGCInterval() time.Duration {
	if v, err := DefaultInt("pool.gc.interval"); err == nil {
		return time.Duration(v) * time.Second
	} else {
		panic(fmt.Errorf("pool.gc.interval Error"))
	}
}

func PoolGCTimeout() int {
	if v, err := DefaultInt("pool.gc.timeout"); err == nil {
		return v
	} else {
		panic(fmt.Errorf("pool.gc.timeout Error"))
	}
}

func ClientRPCDialTimeout(name, version string) time.Duration {
	if v, err := Int(name, version, "client.rpc.dial.timeout"); err == nil {
		return time.Duration(v) * time.Second
	} else {
		if v, err := DefaultInt("client.rpc.dial.timeout"); err == nil {
			return time.Duration(v) * time.Second
		} else {
			panic(fmt.Errorf("client.rpc.dial.timeout Error:%v, %v", name, version))
		}
	}
}

func ClientRPCCallTimeout(name, version string) time.Duration {
	if v, err := Int(name, version, "client.rpc.call.timeout"); err == nil {
		return time.Duration(v) * time.Second
	} else {
		if v, err := DefaultInt("client.rpc.call.timeout"); err == nil {
			return time.Duration(v) * time.Second
		} else {
			panic(fmt.Errorf("client.rpc.call.timeout Error:%v, %v", name, version))
		}
	}
}

func ClientRPCRetry() int {
	if v, err := DefaultInt("client.rpc.retry"); err == nil {
		return v
	} else {
		panic(fmt.Errorf("client.rpc.retry Error"))
	}
}

func ClientConnMax(name, version string) int {
	if v, err := Int(name, version, "client.conn.max"); err == nil {
		return v
	} else {
		if v, err := DefaultInt("client.conn.max"); err == nil {
			return v
		} else {
			panic(fmt.Errorf("client.conn.max Error:%v, %v", name, version))
		}
	}
}

func ClientConnIdle(name, version string) int {
	if v, err := Int(name, version, "client.conn.idle"); err == nil {
		return v
	} else {
		if v, err := DefaultInt("client.conn.idle"); err == nil {
			return v
		} else {
			panic(fmt.Errorf("client.conn.idle Error:%v, %v", name, version))
		}
	}
}

func ServiceCronRegister(name, version string) time.Duration {
	if v, err := Int(name, version, "service.cron.register"); err == nil {
		return time.Duration(v) * time.Second
	} else {
		if v, err := DefaultInt("service.cron.register"); err == nil {
			return time.Duration(v) * time.Second
		} else {
			panic(fmt.Errorf("service.cron.register Error:%v, %v", name, version))
		}
	}
}

package conn

import (
	"fmt"
	"github.com/lemontree2015/skynet"
	"github.com/lemontree2015/skynet/config"
	"github.com/lemontree2015/skynet/rpc/bsonrpc"
	"net"
	"net/rpc"
	"time"
)

// 表示一条TCP连接(连接到RPC Server)
//
// 注意:
// 每次RPC调用都会启动一个goroutine来执行

type Connection struct {
	serviceInfo    *skynet.ServiceInfo  // 连接到的目标ServiceInfo
	conn           net.Conn             // TCP Connect
	clientUUID     string               // handshake之后, RPC Server下发的clientUUID
	rpcClientCodec *bsonrpc.ClientCodec // RPC Codec
	rpcClient      *rpc.Client          // handshake成功后创建RPC Client
	closed         bool                 // 是否关闭
}

func NewConnection(si *skynet.ServiceInfo, network string, timeout time.Duration) (*Connection, error) {
	c, err := net.DialTimeout(network, si.Addr.String(), timeout)

	if err != nil {
		return nil, err
	}

	return NewConnectionFromNetConn(si, c)
}

func NewConnectionFromNetConn(si *skynet.ServiceInfo, c net.Conn) (*Connection, error) {
	cn := &Connection{
		serviceInfo:    si,
		conn:           c,
		clientUUID:     "",
		rpcClientCodec: bsonrpc.NewClientCodec(c),
		rpcClient:      nil,
		closed:         false,
	}

	// 先执行HandShake(成功会创建RPC Client)
	if err := cn.performHandshake(); err == nil {
		return cn, nil
	} else {
		return nil, err
	}
}

func (c *Connection) Close() {
	c.closed = true
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
}

func (c *Connection) IsClosed() bool {
	return c.closed
}

func (c *Connection) ServiceInfo() *skynet.ServiceInfo {
	return c.serviceInfo
}

// Default Timeout
func (c *Connection) Send(fn string, in interface{}, out interface{}) error {
	return c.SendTimeout(fn, in, out, 0)
}

// 发送一次RPC请求
//
// 注意:
// 1. 参数说明
// in  - 客户端请求参数(BSON Encode前的数据)
// out - 服务器的返回结果(BSON Decode后的数据)
//
// 2. 每一次调用都会启动一个独立的goroutine, 结果通过chan返回
func (c *Connection) SendTimeout(fn string, in interface{}, out interface{}, timeout time.Duration) error {
	// 连接已经关闭
	if c.IsClosed() {
		return fmt.Errorf("Connection Closed Error")
	}

	rpcErrChan := make(chan error)

	// 每一次RPC调用都会启一个独立的goroutine
	go func() {
		if c.rpcClient != nil {
			rpcErr := c.rpcClient.Call(c.serviceInfo.Name+"."+fn, in, out) // 进行RPC Call ServiceName.Method
			rpcErrChan <- rpcErr
		}
	}()

	var rpcErr error

	if timeout == 0 {
		timeout = config.ClientRPCCallTimeout(c.serviceInfo.Name, c.serviceInfo.Version)
	}

	t := time.After(timeout)

	select {
	case rpcErr = <-rpcErrChan:
		// RPC Error(直接关闭conn)
		if rpcErr != nil {
			c.Close()
			return rpcErr
		}
	case <-t:
		// RPC Timeout(直接关闭conn)
		c.Close()
		return fmt.Errorf("Connection: timing out request after %s", timeout.String())
	}

	return nil
}

func (c *Connection) performHandshake() error {
	// 读取ServiceRPCServerHandshake
	serviceRPCServerHandshake := &skynet.ServiceRPCServerHandshake{}
	err := c.rpcClientCodec.Decoder.Decode(&serviceRPCServerHandshake)
	if err != nil {
		c.Close()
		return fmt.Errorf("performHandshake Error")
	}

	// 检测 & 赋值
	c.clientUUID = serviceRPCServerHandshake.ClientUUID
	if serviceRPCServerHandshake.ServiceName != c.serviceInfo.Name {
		c.Close()
		return fmt.Errorf("performHandshake Error")
	}

	// 发送ServiceRPCClientHandshake
	serviceRPCClientHandshake := &skynet.ServiceRPCClientHandshake{
		ClientUUID: c.clientUUID,
	}

	err = c.rpcClientCodec.Encoder.Encode(serviceRPCClientHandshake) // 发送ClientHandshake
	if err != nil {
		c.Close()
		return fmt.Errorf("performHandshake Error")
	}

	// 构造RPC Client
	c.rpcClient = rpc.NewClientWithCodec(c.rpcClientCodec)

	return nil
}

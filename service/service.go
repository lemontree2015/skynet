package service

import (
	"github.com/golang/glog"
	"github.com/lemontree2015/skynet"
	"github.com/lemontree2015/skynet/client"
	"github.com/lemontree2015/skynet/config"
	"github.com/lemontree2015/skynet/cron"
	"github.com/lemontree2015/skynet/rpc/bsonrpc"
	"net"
	"net/rpc"
	"runtime/debug"
	"sync"
)

// 核心有2个部分:
// 1. *skynet.ServiceInfo标识一个Service
// 2. delegate表示一个自定义的RPC Service
//
// 会启动2 + N 个goruotine
// 1. 一个负责Accept Loop
// 2. 一个负责处理新进conn以及shutdown等事件
// 以后, 每接收到一个conn, 启动一个goroutine

type Service struct {
	*skynet.ServiceInfo // 每个Service包含一个ServiceInfo

	delegate    interface{}      // 自定义的Skynet RPC Service
	rpcServ     *rpc.Server      // RPC Server
	rpcListener *net.TCPListener // TCP Listener

	clientInfoManager *clientInfoManager
	cronEvery         *cron.CronEvery

	connectionChan chan *net.TCPConn // 每接收到一个新的连接, 通过这个chan发送给mux() loop

	doneChan  chan bool       // 发送信号break mux() loop
	doneGroup *sync.WaitGroup // 整个Service是否退出

	shutdownChan chan bool // shutdown chan
	shuttingDown bool      // Service是否正在shutting, 正在shutting将不再接收新的RPC Request
}

func CreateService(delegate interface{}, si *skynet.ServiceInfo) *Service {
	service := &Service{
		ServiceInfo:       si,
		delegate:          delegate,
		rpcServ:           nil,
		rpcListener:       nil,
		clientInfoManager: NewclientInfoManager(),
		connectionChan:    make(chan *net.TCPConn),
		doneChan:          make(chan bool, 1),
		doneGroup:         new(sync.WaitGroup),
		shutdownChan:      make(chan bool),
		shuttingDown:      false,
	}

	service.cronEvery = cron.NewCronEvery(config.ServiceCronRegister(si.Name, si.Version), service.cronEveryFun) // 注册CronEvery Function

	// 构造RPC Server &  注册RPC
	service.rpcServ = rpc.NewServer()
	service.rpcServ.RegisterName(si.Name, delegate)

	return service
}

func (service *Service) Start() *sync.WaitGroup {
	bindWait := new(sync.WaitGroup)
	bindWait.Add(1)
	go service.tryListen(bindWait) // 启动一个独立的goroutine, Accept TCP Connection
	// 等待Bind & Listen成功
	bindWait.Wait()

	// 启动goroutine
	service.doneGroup.Add(1)
	go func() {
		service.mux()
		service.doneGroup.Done()
	}()

	// Register
	go client.RegisterServices(service.ServiceInfo)

	return service.doneGroup
}

func (service *Service) tryListen(bindWait *sync.WaitGroup) {
	// Bind & Listen
	var err error
	service.rpcListener, err = service.Addr.Listen()
	if err != nil {
		panic(err)
	}
	bindWait.Done()

	// Accept Loop
	for {
		conn, err := service.rpcListener.AcceptTCP()

		if service.shuttingDown {
			break // break loop
		}

		if err != nil && !service.shuttingDown {
			glog.Warningf("AcceptTCP failed: %v", err)
			continue
		}
		service.connectionChan <- conn
	}
}

func (service *Service) Shutdown() {
	if service.shuttingDown {
		return
	}
	service.shutdownChan <- true
}

func (service *Service) cronEveryFun() {
	if service.shuttingDown {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			glog.Infof("Service.cronEveryFun: %v, DEBUG.STACK=%v", r, string(debug.Stack()))
			return
		}
	}()

	// 定期上报自己的信息
	client.RegisterServices(service.ServiceInfo)
}

// Message Loop
func (service *Service) mux() {
loop:
	for {
		select {
		case conn := <-service.connectionChan:
			// 每一个连接启动一个goroutine
			go func() {
				// 执行handshake的逻辑

				// 生成客户端信息
				clientUUID := config.NewUUID()
				clientInfo := &clientInfo{
					ClientUUID: clientUUID,
					Address:    conn.RemoteAddr(),
				}

				// 保存客户端信息
				service.clientInfoManager.Set(clientInfo)

				// 发送ServiceRPCServerHandshake给客户端
				serviceRPCServerHandshake := &skynet.ServiceRPCServerHandshake{
					ClientUUID:  clientUUID,
					ServiceName: service.Name,
				}

				codec := bsonrpc.NewServerCodec(conn) // Codec
				err := codec.Encoder.Encode(serviceRPCServerHandshake)
				if err != nil {
					conn.Close()
					return
				}

				// 读取客户端的ServiceRPCClientHandshake信息
				serviceRPCClientHandshake := &skynet.ServiceRPCClientHandshake{}
				err = codec.Decoder.Decode(serviceRPCClientHandshake)
				if err != nil {
					conn.Close()
					return
				}

				// handshake成功, 启动RPC Service
				service.rpcServ.ServeCodec(codec) // 阻塞当前的goroutine
			}()
		case <-service.shutdownChan:
			service.shutdown()
		case _ = <-service.doneChan:
			break loop
		}
	}
}

// 备注:
// 在RPC goroutine中执行
func (service *Service) shutdown() {
	if service.shuttingDown {
		return
	}

	service.shuttingDown = true

	service.doneGroup.Add(1)
	service.rpcListener.Close()

	service.doneChan <- true // break mux loop

	// UnRegister
	client.UnRegisterServices(service.ServiceInfo)

	service.doneGroup.Done() // 整个Service退出
}

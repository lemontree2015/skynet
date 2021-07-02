package bsonrpc

import (
	"io"
	"net/rpc"
	"testing"
)

// io.ReadWriteCloser
type duplex struct {
	io.Reader
	io.Writer
}

func (d *duplex) Close() error {
	return nil
}

type TestParam struct {
	Val1 string
	Val2 int
	Val3 []string
}

type RPCTest struct {
}

func (ts *RPCTest) Foo(in *TestParam, out *TestParam) error {
	out.Val1 = in.Val1 + "world!"
	out.Val2 = in.Val2 + 5

	//out.Val3 = make([]string, 0, 0)
	for _, v := range in.Val3 {
		out.Val3 = append(out.Val3, v)
	}
	return nil
}

func TestBasicClientServer(t *testing.T) {
	toServer, fromClient := io.Pipe()
	toClient, fromServer := io.Pipe()

	// New RPC Server
	rpcServer := rpc.NewServer()
	ts := &RPCTest{}
	rpcServer.Register(ts) // 注册一个RPC函数
	go rpcServer.ServeCodec(NewServerCodec(&duplex{toServer, fromServer}))

	// New RPC Client
	rpcClient := NewClient(&duplex{toClient, fromClient})

	tp1 := &TestParam{ // Request
		Val1: "Hello ",
		Val2: 4,
		Val3: []string{"1", "2", "3"},
	}
	tp2 := &TestParam{} // Response

	// RPC Call
	err := rpcClient.Call("RPCTest.Foo", tp1, &tp2) // 进行RPC调用
	if err != nil {
		t.Error(err)
		return
	}

	if tp2.Val1 != "Hello world!" {
		t.Errorf("tp2.Val1: Error %d", tp2.Val1)
	}
	if tp2.Val2 != 9 {
		t.Errorf("tp2.Val2: Error %d", tp2.Val2)
	}
	if len(tp2.Val3) != 3 {
		t.Errorf("tp2.Val3: Error %d", tp2.Val3)
	}
	if tp2.Val3[0] != "1" {
		t.Errorf("tp2.Val3: Error %d", tp2.Val3)
	}
	if tp2.Val3[1] != "2" {
		t.Errorf("tp2.Val3: Error %d", tp2.Val3)
	}
	if tp2.Val3[2] != "3" {
		t.Errorf("tp2.Val3: Error %d", tp2.Val3)
	}
}

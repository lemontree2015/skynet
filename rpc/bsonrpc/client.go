package bsonrpc

import (
	"io"
	"net/rpc"
)

type ClientCodec struct {
	conn    io.ReadWriteCloser
	Encoder *Encoder
	Decoder *Decoder
}

func NewClientCodec(conn io.ReadWriteCloser) *ClientCodec {
	return &ClientCodec{
		conn:    conn,
		Encoder: NewEncoder(conn),
		Decoder: NewDecoder(conn),
	}
}

// 1. 写入rpc.Request
// 2. 写入BSON编码的v
func (client *ClientCodec) WriteRequest(req *rpc.Request, v interface{}) error {
	// 1. 写入rpc.Request
	if err := client.Encoder.Encode(req); err != nil {
		client.Close()
		return err
	}

	// 2. 写入BSON编码的v
	if err := client.Encoder.Encode(v); err != nil {
		client.Close()
		return err
	}

	return nil
}

// 1. 读取并解码rpc.Response
func (client *ClientCodec) ReadResponseHeader(res *rpc.Response) error {
	// 1. 读取并解码rpc.Response
	if err := client.Decoder.Decode(res); err != nil {
		client.Close()
		return err
	}

	return nil
}

// 1. 读取并解码V
func (client *ClientCodec) ReadResponseBody(v interface{}) error {
	// 1. 读取并解码V
	if err := client.Decoder.Decode(v); err != nil {
		client.Close()
		return err
	}

	return nil
}

func (client *ClientCodec) Close() (err error) {
	return client.conn.Close()
}

func NewClient(conn io.ReadWriteCloser) (c *rpc.Client) {
	cc := NewClientCodec(conn)
	c = rpc.NewClientWithCodec(cc)
	return
}

package bsonrpc

import (
	"io"
	"net/rpc"
)

type ServerCodec struct {
	conn    io.ReadWriteCloser
	Encoder *Encoder
	Decoder *Decoder
}

func NewServerCodec(conn io.ReadWriteCloser) *ServerCodec {
	return &ServerCodec{
		conn:    conn,
		Encoder: NewEncoder(conn),
		Decoder: NewDecoder(conn),
	}
}

// 1. 读取并解码rpc.Request
func (server *ServerCodec) ReadRequestHeader(rq *rpc.Request) error {
	// 1. 读取并解码rpc.Request
	if err := server.Decoder.Decode(rq); err != nil {
		server.Close()
		return err
	}

	return nil
}

// 1. 读取并解码V
func (server *ServerCodec) ReadRequestBody(v interface{}) error {
	// 1. 读取并解码V
	if err := server.Decoder.Decode(v); err != nil {
		server.Close()
		return err
	}

	return nil
}

// 1. 写入rpc.Response
// 2. 写入BSON编码的v
func (server *ServerCodec) WriteResponse(res *rpc.Response, v interface{}) error {
	// 1. 写入rpc.Response
	if err := server.Encoder.Encode(res); err != nil {
		server.Close()
		return err
	}

	// 2. 写入BSON编码的v
	if err := server.Encoder.Encode(v); err != nil {
		server.Close()
		return err
	}

	return nil
}

func (server *ServerCodec) Close() (err error) {
	return server.conn.Close()
}

func ServeConn(conn io.ReadWriteCloser) (s *rpc.Server) {
	s = rpc.NewServer()
	s.ServeCodec(NewServerCodec(conn))
	return
}

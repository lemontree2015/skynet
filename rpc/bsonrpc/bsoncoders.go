package bsonrpc

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io"
)

// Encoder & Decoder

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// 编码v, 将序列化的数据写到io.Writer
func (encoder *Encoder) Encode(v interface{}) error {
	if v == nil {
		panic(fmt.Errorf("Encode V is Nil"))
	}

	// 编码v
	buf, err := bson.Marshal(v)
	if err != nil {
		return err
	}

	// Send编码后的V
	n, err := encoder.w.Write(buf)
	if err != nil {
		return err
	}

	// 检测发送的长度
	if l := len(buf); n != l {
		return fmt.Errorf("Encode Send Error: Wrote %v bytes, should have wrote %v", n, l)
	}

	return nil
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// 从io.Reader读取数据, 并解码到pv
//
// 备注:
// bson.Marshal编码会自动写入4 bytes的长度作为Header
// | 4 bytes Header | Payload | 4 bytes Header | Payload |....
func (decoder *Decoder) Decode(pv interface{}) error {
	if pv == nil {
		panic(fmt.Errorf("Decode PV is Nil"))
	}

	// 读取4 bytes header
	var lenBuf [4]byte
	n, err := decoder.r.Read(lenBuf[:])

	if n != 4 {
		return fmt.Errorf("Corrupted BSON stream: could only read %v", n)
	}

	if err != nil {
		return err
	}

	// 计算4 bytes header的length
	length := (int(lenBuf[0]) << 0) |
		(int(lenBuf[1]) << 8) |
		(int(lenBuf[2]) << 16) |
		(int(lenBuf[3]) << 24)

	buf := make([]byte, length)
	copy(buf[0:4], lenBuf[:])

	n, err = io.ReadFull(decoder.r, buf[4:]) // 读满buf

	if err != nil {
		return err
	}

	// 检测读取的长度
	if n+4 != length {
		return fmt.Errorf("Decode Error: Expected %v bytes, read %v", length, n)
	}

	return bson.Unmarshal(buf, pv)
}

package bsonrpc

import (
	"bytes"
	"gopkg.in/mgo.v2/bson"
	"net/rpc"
	"testing"
)

func TestEncode(t *testing.T) {
}

func TestDecode(t *testing.T) {
	// Encode
	req := &rpc.Request{
		ServiceMethod: "Foo.Bar",
		Seq:           3,
	}

	b, err := bson.Marshal(req)

	if err != nil {
		t.Fatal(err)
	}

	// Decode
	buf := bytes.NewBuffer(b)
	dec := NewDecoder(buf)

	r := &rpc.Request{}
	err = dec.Decode(r)

	if err != nil {
		t.Fatal(err)
	}

	if req.ServiceMethod != r.ServiceMethod ||
		req.Seq != r.Seq {
		t.Fatal("Values don't match")
	}
}

func TestDecodeReadsOnlyOne(t *testing.T) {

	// Encode 2 values
	req := &rpc.Request{
		ServiceMethod: "Foo.Bar",
		Seq:           3,
	}

	type T struct {
		Value string
	}

	tv := &T{
		Value: "test",
	}

	b, err := bson.Marshal(req)

	if err != nil {
		t.Fatal(err)
	}

	buf := bytes.NewBuffer(b)

	b1, err := bson.Marshal(tv)

	if err != nil {
		t.Fatal(err)
	}

	buf.Write(b1) // 继续写入tv

	// Decode 2 values
	dec := NewDecoder(buf)

	r := &rpc.Request{}
	err = dec.Decode(r)

	if req.ServiceMethod != r.ServiceMethod ||
		req.Seq != r.Seq {
		t.Fatal("Values don't match")
	}

	if err != nil {
		t.Fatal(err)
	}

	// We should be able to read a second message off this io.Reader
	tmp := &T{}
	err = dec.Decode(tmp)

	if err != nil {
		t.Fatal(err)
	}

	if tmp.Value != tv.Value {
		t.Fatal("Values don't match")
	}

}

package config

import (
	"crypto/rand"
	"fmt"
	"io"
)

// UUID生成器
//
// e.g.
// 33f2fa77-394b-4a06-93f3-0e7902353d85
// d618c688-7e0c-4f4e-9691-d828992c456f
// a9766ca8-db41-4828-8845-82378419f3c7
func NewUUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(fmt.Errorf("NewUUID Error:%v", err))
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

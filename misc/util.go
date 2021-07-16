package misc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/lemontree2015/skynet/logger"
	"hash/crc64"
	"os"
	"time"
)

func CheckFatal(err error) {
	if err != nil {
		logger.Logger.Fatalf("CheckFatal: %v", err)
		os.Exit(1)
	}
}

// 返回和下面MySQL函数等价的结果
// SELECT UNIX_TIMESTAMP();
func UnixTimestamp() int64 {
	return time.Now().Unix()
}

func UnixNanoTimestamp() int64 {
	return time.Now().UnixNano()
}

func RandomSeed() int64 {
	return time.Now().UnixNano()
}

// 返回长度为32的MD5 Str(小写形式)
func MD5Bytes(data []byte) string {
	m := md5.Sum(data)
	return hex.EncodeToString(m[:])
}

// 返回长度为32的MD5 Str(小写形式)
func MD5Str(str string) string {
	m := md5.Sum([]byte(str))
	return hex.EncodeToString(m[:])
}

func Hash(i string) string {
	s := md5.Sum([]byte(i))
	return fmt.Sprintf("0x%x", crc64.Checksum(s[:], crc64.MakeTable(crc64.ISO)))
}

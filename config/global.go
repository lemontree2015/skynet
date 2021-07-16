package config

import (
	"flag"
	"fmt"
)

// 全局唯一的配置文件

var (
	skynetCFGPath      string // --conf_path=xxx
	skynetInstanceUUID string // --instance_uuid=xxx
)

// 全局配置信息
var (
	skynetCFG *Configuration
)

// 解析配置文件
func init() {
	flag.StringVar(&skynetCFGPath, "conf_path", "", "skynet configuration file")
	flag.StringVar(&skynetInstanceUUID, "instance_uuid", "", "skynet instance uuid")

	// 注意: 这里一定要调用flag.Parse来解析参数, 否则无法获取*.conf的路径
	flag.Parse()

	if skynetCFGPath == "" {
		panic(fmt.Errorf("skynetCFGPath is Empty"))
	}

	if skynetInstanceUUID == "" {
		skynetInstanceUUID = NewUUID()
	}

	Reload()
}

func Reload() {
	// 解析skynet.conf
	skynetCFGTmp, err := Parse(skynetCFGPath)
	if err != nil {
		panic(fmt.Errorf("Parse Config Error:%v", err))
	}
	skynetCFG = skynetCFGTmp // 替换
	//fmt.Printf("parse skynet configuration success, path=%v", skynetCFG.FilePath())
}

func FilePath() string {
	return skynetCFG.FilePath()
}

func InstanceUUID() string {
	return skynetInstanceUUID
}

func DefaultString(option string) (value string, err error) {
	return skynetCFG.DefaultString(option)
}

func DefaultInt(option string) (value int, err error) {
	return skynetCFG.DefaultInt(option)
}

func DefaultBool(option string) (value bool, err error) {
	return skynetCFG.DefaultBool(option)
}

func String(serviceName, serviceVersion, option string) (value string, err error) {
	return skynetCFG.String(serviceName, serviceVersion, option)
}

func Int(serviceName, serviceVersion, option string) (value int, err error) {
	return skynetCFG.Int(serviceName, serviceVersion, option)
}

func Bool(serviceName, serviceVersion, option string) (value bool, err error) {
	return skynetCFG.Bool(serviceName, serviceVersion, option)
}

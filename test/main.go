package main

import (
	"fmt"
	"monitor_server/skynet"
	"monitor_server/skynet/client"
)

func main() {
	fmt.Println("run test")
	fmt.Println(client.Get(skynet.NewCriteria()))
}

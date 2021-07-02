package main

import (
	"fmt"
	"github.com/lemontree2015/skynet"
	"github.com/lemontree2015/skynet/client"
)

func main() {
	fmt.Println("run test")
	fmt.Println(client.Get(skynet.NewCriteria()))
}

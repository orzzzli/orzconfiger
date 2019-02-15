package main

import (
	"fmt"
	"orzconfiger"
	"time"
)

func main() {
	orzconfiger.InitConfiger("")

	for {
		fmt.Println(orzconfiger.ConfigerMap)
		time.Sleep(time.Duration(10) * time.Second)
	}
}

package main

import (
	"fmt"
	"orzconfiger"
	"time"
)

func main() {
	orzconfiger.InitConfiger("")

	for {
		fmt.Println(orzconfiger.ConfigerMap,orzconfiger.ConfigerSection)
		time.Sleep(time.Duration(10) * time.Second)
	}
}

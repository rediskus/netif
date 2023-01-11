package main

import (
	"fmt"
	ni "github.com/rediskus/netif"
)

func main() {
	is, _ := ni.Parse(
		ni.Path("./input"),
	)
	if e := is.FindAdapter("eth0"); e != nil {
		str := e.DNSConcatString()
		fmt.Println(str)
		arr := e.DNSStrings()
		fmt.Println(arr)
	}
	_ = is.Write(
		"eth0", ni.Path("./output"),
	)
}

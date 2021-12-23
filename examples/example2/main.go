package main

import (
	"fmt"
	ni "github.com/rediskus/netif"
)

func main() {
	i, _ := ni.Parse(ni.Path("./interfaces"))
	for _, r := range i.Adapters {
		fmt.Printf("Name %s\r\n", r.Name)
		fmt.Printf("Network %s\r\n", r.Network)
		fmt.Printf("Address %s\r\n", r.Address)
		fmt.Printf("Netmask %s\r\n", r.Netmask)
		fmt.Printf("Gateway %s\r\n", r.Gateway)
	}
	i.Write("eth1", ni.Path("./interfaces1"))
}

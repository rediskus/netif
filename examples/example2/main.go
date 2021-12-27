package main

import (
	"fmt"
	ni "github.com/rediskus/netif"
	"net"
)

func main() {
	i, _ := ni.Parse(ni.Path("./interfaces"))
	eth := i.FindAdapter("eth1")
	if eth != nil {
		fmt.Printf("Adapter found. Deleting...")
		if i.DeleteAdapter(eth.Name) == true {
			fmt.Println("done")
		} else {
			fmt.Println("error")
		}
	}
	if e, err := i.AddAdapter("eth2"); err == nil {
		e.AddrSource = ni.STATIC
		e.Address = net.ParseIP("192.168.1.10")
		e.Netmask = net.ParseIP("255.255.255.0")
		e.Gateway = net.ParseIP("192.168.1.1")
		e.SetPreUp("wb-gsm restart_if_broken")
		e.SetPreUp("sleep 10")
		e.SetPreUp("wb-gsm off")
		e.Hotplug = true
	} else {
		fmt.Println(err.Error())
	}
	_ = i.Write("", ni.Path("./interfaces1"))
}

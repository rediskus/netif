package main

import ni "github.com/swha12/netif"

func main() {
	is, _ := ni.Parse(
		ni.Path("input"),
	)

	is.Write(
		"eth0", ni.Path("output"),
	)
}

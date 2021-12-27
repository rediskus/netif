package netif

import (
	"errors"
	"fmt"
	"net"
)

type AddrSource int

const (
	DHCP AddrSource = 1 + iota
	STATIC
	LOOPBACK
	MANUAL
)

type AddrFamily int

const (
	INET AddrFamily = 1 + iota
	INET6
)

// A representation of a network adapter
type NetworkAdapter struct {
	Name       string
	Hotplug    bool
	Auto       bool
	Address    net.IP
	Netmask    net.IP
	Network    net.IP
	Broadcast  net.IP
	Gateway    net.IP
	AddrSource AddrSource
	AddrFamily AddrFamily
	//[added] DNS namserver added to read/write dns-nameservers
	DNSNS []net.IP
	// dsk
	isBridge       bool
	BridgePorts    []string
	BridgeWaitport string
	BridgeFd       string
	BridgeMaxwait  string
	BridgeStp      bool
	PreUp          []string
	Hostname       string
}

func (na *NetworkAdapter) validateIP(strIP string) (net.IP, error) {
	var ip net.IP
	if ip = net.ParseIP(strIP); ip == nil {
		return nil, errors.New("invalid IP address")
	}
	return ip, nil
}

func (na *NetworkAdapter) SetAddress(address string) error {
	addr, err := na.validateIP(address)
	if err != nil {
		return err
	}
	na.Address = addr
	return nil
}

func (na *NetworkAdapter) SetNetmask(address string) error {
	addr, err := na.validateIP(address)
	if err == nil {
		na.Netmask = addr
	}
	return err
}

func (na *NetworkAdapter) SetGateway(address string) error {
	addr, err := na.validateIP(address)
	if err == nil {
		na.Gateway = addr
	}
	return err
}

func (na *NetworkAdapter) SetBroadcast(address string) error {
	addr, err := na.validateIP(address)
	if err == nil {
		na.Broadcast = addr
	}
	return err
}

func (na *NetworkAdapter) SetNetwork(address string) error {
	addr, err := na.validateIP(address)
	if err == nil {
		na.Network = addr
	}
	return err
}

//SetDNSNS sets DNS Nameservers
//[added]
func (na *NetworkAdapter) SetDNSNS(address string) error {
	addr, err := na.validateIP(address)
	if err == nil {
		na.DNSNS = append(na.DNSNS, addr)
	}
	return err
}

func (na *NetworkAdapter) SetConfigType(configType string) error {
	switch configType {
	case "DHCP":
		na.AddrSource = DHCP
	case "STATIC":
		na.AddrSource = STATIC
	default:
		return fmt.Errorf("unexpected configType: %s", configType)
	}
	return nil
}

func (na *NetworkAdapter) SetBridgePort(port string) {
	na.BridgePorts = append(na.BridgePorts, port)
	na.isBridge = true
}

func (na *NetworkAdapter) ParseAddressSource(AddressSource string) (AddrSource, error) {
	// Parse the address source for an interface
	var src AddrSource
	switch AddressSource {
	case "static":
		src = STATIC
	case "dhcp":
		src = DHCP
	case "loopback":
		src = LOOPBACK
	case "manual":
		src = MANUAL
	default:
		return -1, errors.New("invalid address source")
	}
	return src, nil
}

func (na *NetworkAdapter) ParseAddressFamily(AddressFamily string) (AddrFamily, error) {
	// Parse the address family for an interface
	var fam AddrFamily
	switch AddressFamily {
	case "inet":
		fam = INET
	case "inet6":
		fam = INET6
	default:
		return -1, errors.New("invalid address family")

	}
	return fam, nil
}

func (na *NetworkAdapter) DNSConcatString() string {
	message := ""
	for indx, dns := range na.DNSNS {
		if indx > 0 || indx < len(na.DNSNS)-1 {
			message = message + "," + dns.String()
		} else {
			message = message + dns.String()
		}
	}
	return message
}

func (na *NetworkAdapter) GetAddrFamilyString() string {
	switch na.AddrFamily {
	case INET:
		return "inet"
	case INET6:
		return "inet6"
	}
	return "inet"
}

func (na *NetworkAdapter) GetSourceFamilyString() string {
	switch na.AddrSource {
	case DHCP:
		return "dhcp"
	case STATIC:
		return "static"
	case LOOPBACK:
		return "loopback"
	case MANUAL:
		return "manual"
	}
	return "dhcp"
}

func ip2string(ip net.IP) string {
	if len(ip)>0 {
		return ip.String()
	}
	return ""
}

func (na *NetworkAdapter) GetAddress() string {
	return ip2string(na.Address)
}

func (na *NetworkAdapter) GetNetmask() string {
	return ip2string(na.Netmask)
}

func (na *NetworkAdapter) GetGateway() string {
	return ip2string(na.Gateway)
}
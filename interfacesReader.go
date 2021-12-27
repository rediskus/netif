package netif

import (
	"bufio"
	"os"
	"strings"

"github.com/n-marshall/fn"
)

// TODO get rid of interfaceReader
type InterfacesReader struct {
	filePath    string
	adapters    []*NetworkAdapter
	autoList    []string
	hotplugList []string
	context     int
}

func Parse(opts ...fn.Option) (*InterfaceSet, error) {
	var err error

	fnConfig := fn.MakeConfig(
		fn.Defaults{"path": "/etc/network/interfaces"},
		opts,
	)
	path := fnConfig.GetString("path")

	is := &InterfaceSet{
		InterfacesPath: path,
	}
	if is.Adapters, err = NewInterfacesReader(is.InterfacesPath).ParseInterfaces(); err != nil {
		return nil, err
	}

	return is, nil
}

func NewInterfacesReader(filePath string) *InterfacesReader {
	ir := InterfacesReader{filePath: filePath}
	ir.reset()

	return &ir
}

func (ir *InterfacesReader) ParseInterfaces() ([]*NetworkAdapter, error) {
	// Reset this object in case is not new
	ir.reset()

	// Try to open the file
	f, err := os.Open(ir.filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Treat each line from the file
	ir.readLinesFromFile(f)

	return ir.parseInterfacesImplementation(), nil
}

func (ir *InterfacesReader) parseInterfacesFromString(data string) {
	// Reset this object in case is not new
	ir.reset()
}

func (ir *InterfacesReader) parseInterfacesImplementation() []*NetworkAdapter {
	// Save adapters and return them

	// foreach iface in the auto list
	for _, autoName := range ir.autoList {
		for naIdx, _ := range ir.adapters {
			if ir.adapters[naIdx].Name == autoName {
				ir.adapters[naIdx].Auto = true
			}
		}
	}

	// foreach iface in the hotplug list
	for _, hotplugName := range ir.hotplugList {
		for naIdx, _ := range ir.adapters {
			if ir.adapters[naIdx].Name == hotplugName {
				ir.adapters[naIdx].Hotplug = true
			}
		}
	}

	return ir.adapters
}

func (ir *InterfacesReader) readLinesFromFile(file *os.File) bool {
	s := bufio.NewScanner(file)

	//var a Adapter

	for s.Scan() {
		line := s.Text()

		// Identify the clauses by analyzing the first word of each line.
		// Go to the next line if the current line is a comment.
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		// Continue if line is empty
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// Parse the line
		ir.parseIface(line)
		ir.parseDetails(line)
		ir.readAuto(line)
		ir.readHotplug(line)
	}
	return false
}

func (ir *InterfacesReader) parseIface(line string) {
	if !strings.HasPrefix(line, "iface") {
		return
	}

	sline := strings.Split(strings.TrimSpace(line), " ")

	ir.adapters = append(ir.adapters, &NetworkAdapter{Name: sline[1]})
	ir.context++

	// Parse and set the address source
	src, err := ir.adapters[ir.context].ParseAddressSource(sline[len(sline)-1])
	if err != nil {
		panic(err)
	}
	ir.adapters[ir.context].AddrSource = src

	// Parse and set the address family
	fam, err := ir.adapters[ir.context].ParseAddressFamily(sline[2])
	if err == nil {
		ir.adapters[ir.context].AddrFamily = fam
	}
}

func (ir *InterfacesReader) parseDetails(line string) {
	// If line begins with a space, it's a interface attribute
	if strings.TrimSpace(line)[0] == line[0] {
		// Doesn't begin with space, pass
		isAttr := false
		sline := strings.Split(strings.TrimSpace(line), " ")
		switch sline[0] {
		case "address":
			isAttr = true
		case "netmask":
			isAttr = true
		case "network":
			isAttr = true
		case "broadcast":
			isAttr = true
		case "gateway":
			isAttr = true
		case "dns-nameservers":
			isAttr = true
		case "dns-domain":
			isAttr = true
		case "dns-search":
			isAttr = true
		case "bridge_ports":
			isAttr = true
		case "bridge_waitport":
			isAttr = true
		case "bridge_stp":
			isAttr = true
		case "bridge_fd":
			isAttr = true
		case "bridge_maxwait":
			isAttr = true
		case "pre-up":
			isAttr = true
		case "hostname":
			isAttr = true
		}
		if isAttr == false {
			return
		}
	}

	sline := strings.Split(strings.TrimSpace(line), " ")
	na := ir.adapters[ir.context]

	switch sline[0] {
	case "address":
		if na.SetAddress(sline[1]) != nil {
			return
		}
	case "netmask":
		if na.SetNetmask(sline[1]) != nil {
			return
		}
	case "gateway":
		if na.SetGateway(sline[1]) != nil {
			return
		}
	case "broadcast":
		if na.SetBroadcast(sline[1]) != nil {
			return
		}
	case "network":
		if na.SetNetwork(sline[1]) != nil {
			return
		}
	//[added] read dns-nameservers list
	case "dns-nameservers":
		//na.DNSNS = make([]net.IP, 0, 4)
		for i := 1; i < len(sline); i++ {
			if na.SetDNSNS(sline[i]) != nil {
				continue
			}
		}
	// dsk
	case "bridge_ports":
		for i:=1;i<len(sline);i++ {
			na.BridgePorts = append(na.BridgePorts, sline[i])
		}
		na.isBridge = true
	case "bridge_waitport":
		na.BridgeWaitport = sline[1]
	case "bridge_stp":
		if strings.ToUpper(sline[1]) == "OFF" {
			na.BridgeStp = false
		} else {
			na.BridgeStp = true
		}
	case "bridge_fd":
		na.BridgeFd = sline[1]
	case "bridge_maxwait":
		na.BridgeMaxwait = sline[1]
	case "pre-up":
		if len(na.PreUp)>0 {
			na.PreUp = append(na.PreUp,":") // разделяем разные pre-up
		}
		for i:=1;i<len(sline);i++ {
			na.PreUp = append(na.PreUp, sline[i])
		}
	case "hostname":
		na.Hostname = sline[1]
	default:
	}
}

func (ir *InterfacesReader) readWord(line string, word string) (bool, string) {
	// Isolate the second value after a matching word on the given line
	if strings.HasPrefix(line, word) {
		sline := strings.Split(strings.TrimSpace(line), " ")
		for _, s := range sline {
			if s != word {
				return true, s
			}
		}
	}
	return false, ""
}

func (ir *InterfacesReader) readAuto(line string) {
	// Identify which adapters are flagged auto
	if ok, iface := ir.readWord(line, "auto"); ok {
		ir.autoList = append(ir.autoList, iface)
	}
}

func (ir *InterfacesReader) readHotplug(line string) {
	// Identify which adapters are flagged allow-hotplug
	if ok, iface := ir.readWord(line, "allow-hotplug"); ok {
		ir.hotplugList = append(ir.hotplugList, iface)
	}
}

func (ir *InterfacesReader) reset() {
	// Initialize a place to store create NetworkAdapter objects
	ir.adapters = nil

	// Keep a list of adapters that have the auto or allow-hotplug flags set.
	ir.autoList = nil
	ir.hotplugList = nil

	// Store the interface context.
	// This is the index of the adapters collection.
	ir.context = -1
}

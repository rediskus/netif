package netif

import (
	"fmt"
	"github.com/n-marshall/fn"
	"strings"
)

type InterfaceSet struct {
	InterfacesReader

	InterfacesPath string
	Adapters       []*NetworkAdapter
}

func NewInterfaceSet(opts ...fn.Option) *InterfaceSet {
	fnConfig := fn.MakeConfig(
		fn.Defaults{"path": "/etc/network/interfaces"},
		opts,
	)
	path := fnConfig.GetString("path")

	return &InterfaceSet{
		InterfacesPath: path,
	}
}

func Path(path string) fn.Option {
	return fn.MakeOption("path", path)
}

func (i *InterfaceSet) FindAdapter(name string) *NetworkAdapter {
	name = strings.ToLower(name)
	for _,a:=range i.Adapters {
		if a.Name == name {
			return a
		}
	}
	return nil
}

func (i *InterfaceSet) DeleteAdapter(name string) bool {
	name = strings.ToLower(name)
	for idx,a:=range i.Adapters {
		if a.Name == name {
			i.Adapters = append(i.Adapters[:idx],i.Adapters[idx+1:]...)
			return true
		}
	}
	return false
}

func (i *InterfaceSet) AddAdapter(name string) (*NetworkAdapter, error) {
	if i.FindAdapter(name) != nil {
		return nil, fmt.Errorf("Adapter %s already exists",name)
	}

	var eth = NetworkAdapter{Name: name}

	i.Adapters = append(i.Adapters,&eth)

	return &eth,nil
}
package utils

import (
	"errors"
	"regexp"
)

//to succesfully create an inventory plugin:
//  - a variable of any type that implements this interface must be exported
//  - the variable exported above must be named "Inventory"
type Inventory interface {
	FetchInvetory(conf *InventoryConf)
	GetInvetory() []Host
	GetSpec() *InventorySpec
}

//used by inventory plugin to define informations about itself
type InventorySpec struct {
	Name        string
	Description string
	Version     string
	Config      InventoryConf
}

type InventoryConf map[string]string

func GetHost(p Inventory, host string) (*Host, error) {
	re, _ := regexp.Compile(host)
	for _, hst := range p.GetInvetory() {
		if re.MatchString(hst.Fqdn) {
			return &hst, nil
		}
	}
	return nil, errors.New("host not found!")
}

//used by inventory plugin to define host
type Host struct {
	Fqdn        string
	Ipaddresses []string
	Powerstate  string
	Disks       map[string]Disk
}

type Disk struct {
	Path           string
	DiskType       string
	Serial         string
	FirmWare       string
	Manufactor     string
	Backing_device *Disk
	Cache_mode     string
	TotalByte      uint64
	UsedByte       uint64
	TotalInode     uint64
	UsedInode      uint64
}

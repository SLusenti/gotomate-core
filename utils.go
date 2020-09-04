package utils

import "regexp"

//this packageprovide the common struct and interface for build plugins

//to succesfully create a netq-tool plugin:
//  - a variable of any type that implements this interface must be exported
//  - the variable exported above must be named "Pkg"
type Pkg interface {
	Run(subcommand string, mainarg string, config *Config, inv Inventory)
	GetCommand() *Command
}

//option (or flags) definition mapped in each subcommand
//E.G netq-tool <command> <subcommand> [ -l (this is an optional flag) ]
//option can have by parametrized by set Option.HasField = true
//option that has Option.IsRequired = true implicity must have a filed
//field is the default value
//if ConfName is not empty string than an option is written in the conf as below
//cmd:
//  `cmd-name`:
//    `ConfName`: `Field`
type Option struct {
	IsRequired  bool
	HasField    bool
	Field       string
	Description string
	ConfName    string
}

//used by netq-tool plugins to provide a list of actions
type Subcommand struct {
	Options     map[string]*Option
	Description string
}

//used by netq-tool plugins to define command (or plugin name spec)
//Command.command is the "plugin name" used by netq-tool
type Command struct {
	Version     string
	Command     string
	Subcommands map[string]*Subcommand
	Description string
}

//to succesfully create an inventory plugin:
//  - a variable of any type that implements this interface must be exported
//  - the variable exported above must be named "Inventory"
type Inventory interface {
	FetchInvetory()
	GetInvetory() []Host
	GetSpec() *InventorySpec
}

//used by inventory plugin to define host
type Host struct {
	Fqdn        string
	Ipaddresses []string
	Powerstate  string
	Disks       map[string]BlockDisk
}

type BlockDisk interface {
	GetDisk() *Disk
}

type Disk struct {
	Path     string
	DiskType string
}

func (d Disk) GetDisk() *Disk {
	return &d
}

type PhyDisk struct {
	Disk
	Serial string
}

type VirtDisk struct {
	Disk
}

type Bcache struct {
	VirtDisk
	Backing_device BlockDisk
	Cache_mode     string
}

//used by inventory plugin to define informations about itself
type InventorySpec struct {
	Description string
	Version     string
}

//netq-tool configurations (currently config.yaml)
//WIP: not all the parameters are currently used (commented means not used)
//  - JUJU and MAAS config must be in a separate config (like inventory.yaml)
type Config struct {
	SshUser string `yaml:"SshUser"`
	SshKey  string `yaml:"SshKey"`
	Cmd     map[string]map[string]string
	Inv     map[string]map[string]string
}

func GetHost(p *Inventory, host string) (*Host, error) {
	re, _ := regexp.Compile(host)
	for _, hst := range p.GetInvetory() {
		if re.MatchString(hst.Fqdn) {
			return &hst, nil
		}
	}
	return nil, errors.New("host not found!")
}

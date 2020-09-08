package utils

import (
	"errors"
	"regexp"
)

//this packageprovide the common struct and interface for build plugins

//to succesfully create a netq-tool plugin:
//  - a variable of any type that implements this interface must be exported
//  - the variable exported above must be named "Pkg"
type Pkg interface {
	Run(subcommand string, mainarg string, config *PkgConf, inv *Inventory, sshconf *SshConfig)
	GetCommand() *Command
}

//this config is a specific plugin config than will be nested in the conf.yml
//pkgconf:
//  [command]:
//     key: value
type PkgConf map[string]string

//option (or flags) definition mapped in each subcommand
//E.G netq-tool <command> <subcommand> [ -l (this is an optional flag) ] [mainarg]
//option can have by parametrized by set Option.HasField = true
//option that has Option.IsRequired = true implicity must have a filed
//field is the default value
type Option struct {
	IsRequired  bool
	HasField    bool
	FieldName   string
	Values      []string
	Description string
}

//used by netq-tool plugins to provide a list of actions
type Subcommand struct {
	Options     map[string]*Option
	Description string
	HasMainarg  bool
	MainArgDscr string
}

//used by netq-tool plugins to define command (or plugin name spec)
//Command.command is the "plugin name" used by netq-tool
//E.G netq-tool <command> <subcommand> [ -l (this is an optional flag) ] [mainarg]
type Command struct {
	Version     string
	Command     string
	Subcommands map[string]*Subcommand
	Description string
	Config      PkgConf
}

//to succesfully create an inventory plugin:
//  - a variable of any type that implements this interface must be exported
//  - the variable exported above must be named "Inventory"
type Inventory interface {
	FetchInvetory(conf *InventoryConf)
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
	Name        string
	Description string
	Version     string
	Config      InventoryConf
}

type InventoryConf map[string]string

//netq-tool configurations (currently config.yaml)
//WIP: not all the parameters are currently used (commented means not used)
//  - JUJU and MAAS config must be in a separate config (like inventory.yaml)
type Config struct {
	GlobalConf
	SshConfig
	InventoryConf map[string]*InventoryConf
	PkgConf       map[string]*PkgConf
}

type GlobalConf struct {
	LogLevel string
	LogFile  string
}

type SshConfig struct {
	SshUser string
	SshKey  string
}

func GetHost(p Inventory, host string) (*Host, error) {
	re, _ := regexp.Compile(host)
	for _, hst := range p.GetInvetory() {
		if re.MatchString(hst.Fqdn) {
			return &hst, nil
		}
	}
	return nil, errors.New("host not found!")
}

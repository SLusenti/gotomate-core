package utils

//to succesfully create a netq-tool plugin:
//  - a variable of any type that implements this interface must be exported
//  - the variable exported above must be named "Pkg"
type Pkg interface {
	Run(subcommand string, mainarg string, config *PkgConf, inv *Inventory, sshconf *SshConfig)
	GetCommand() *Command
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

//this config is a specific plugin config than will be nested in the conf.yml
//pkgconf:
//  [command]:
//     key: value
type PkgConf map[string]string

//used by netq-tool plugins to provide a list of actions
type Subcommand struct {
	Options     map[string]*Option
	Description string
	HasMainarg  bool
	MainArgDscr string
}

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

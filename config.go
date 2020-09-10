package utils

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

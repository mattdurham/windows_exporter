package collector

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var ConfigMap = make(map[string]func() Config)

func GenerateConfigsWithKingpin(ka *kingpin.Application) map[string]Config {
	configs := make(map[string]Config)
	for k, v := range ConfigMap {
		c := v()
		c.LoadConfigFromKingPin(ka)
		configs[k] = c
	}
	return configs
}

// Used to hold the metadata about a configuration option
type Config interface {
	LoadConfigFromKingPin(ka *kingpin.Application)
	LoadConfigFromMap(m map[string]string)
}

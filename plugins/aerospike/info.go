package aerospike

import (
	"github.com/ehazlett/interlock"
)

const (
	pluginName        = "aerospike"
	pluginVersion     = "0.2"
	pluginDescription = "Aerospike plugin"
	pluginUrl         = "https://github.com/aerospike/interlock/tree/master/plugins/aerospike"
)

var (
	pluginInfo = &interlock.PluginInfo{
		Name:        pluginName,
		Version:     pluginVersion,
		Description: pluginDescription,
		Url:         pluginUrl,
	}
)

package aerospike

import (
	"github.com/ehazlett/interlock"
)

const (
	pluginName        = "aerospike"
	pluginVersion     = "0.1"
	pluginDescription = "AeroSpike plugin"
	pluginUrl         = "https://github.com/rguo-aerospike/interlock/tree/master/plugins/aerospike"
)

var (
	pluginInfo = &interlock.PluginInfo{
		Name:        pluginName,
		Version:     pluginVersion,
		Description: pluginDescription,
		Url:         pluginUrl,
	}
)

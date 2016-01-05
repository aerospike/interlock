package interlock

type (
	Config struct {
		SwarmUrl       string   `json:"swarm_url,omitempty"`
		TLSCaCert		string	`json:"tls_ca_path,omitempty"`
		TLSCert			string	`json:"tls_cert_path,omitempty"`
		TLSKey			string	`json:"tls_key_path,omitempty"`
		EnabledPlugins []string `json:"enabled_plugins,omitempty"`
	}

	InterlockConfig struct {
		Version string `json:"version,omitempty"`
	}
)

package aerospike

type PluginConfig struct {
    ClusterName     string  `json:"cluster_name,omitempty"`
    NetworkName     string  `json:"network_name,omitempty"`
    MeshPort		string	`json:"cluster_mesh_port,omitempty"`

}

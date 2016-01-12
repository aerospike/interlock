package aerospike

import (
	"fmt"
	"time"
    "os"
    "strings"
    "os/exec"
	log "github.com/Sirupsen/logrus"
	"github.com/ehazlett/interlock"
	"github.com/ehazlett/interlock/plugins"
	"github.com/samalba/dockerclient"
)

type AerospikePlugin struct {
    pluginConfig    *PluginConfig
	interlockConfig *interlock.Config
	client          *dockerclient.DockerClient
}

func init() {
	plugins.Register(
		pluginInfo.Name,
		&plugins.RegisteredPlugin{
			New: NewPlugin,
			Info: func() *interlock.PluginInfo {
				return pluginInfo
			},
		})
}

func loadPluginConfig() (*PluginConfig, error) {
    cfg := &PluginConfig{
        ClusterName:    "aerospike",
		NetworkName:	"docker",
		MeshPort:		"3002",
    }
    clusterName  := os.Getenv("AEROSPIKE_CLUSTER_NAME")
    if clusterName != "" {
        cfg.ClusterName = clusterName
    }
    networkName := os.Getenv("AEROSPIKE_NETWORK_NAME")
    if networkName != "" {
        cfg.NetworkName = networkName
    }
    meshPort := os.Getenv("AEROSPIKE_MESH_PORT")
    if meshPort != "" {
        cfg.MeshPort = meshPort
    }
    return cfg, nil
}

func NewPlugin(interlockConfig *interlock.Config, client *dockerclient.DockerClient) (interlock.Plugin, error) {
	//return ExamplePlugin{interlockConfig: interlockConfig, client: client}, nil
    pluginConfig, err := loadPluginConfig()
    if err != nil {
        return nil, err
    }

    plugin := AerospikePlugin{
        pluginConfig:    pluginConfig,
        interlockConfig: interlockConfig,
        client:          client,
    }

    return plugin, nil
}

func (p AerospikePlugin) Info() *interlock.PluginInfo {
	return pluginInfo
}

func (p AerospikePlugin) HandleEvent(event *dockerclient.Event) error {
    switch event.Status {
    case "start":
        if err := p.clusterAerospike(event.Id[:12]); err != nil{
            return err
        }
    }
    return nil
}

func (p AerospikePlugin) contains(containers []string, eventId string) bool {
	for _,a := range containers {
		if a == eventId{
			return true
		}
	}
	return false
}
func (p AerospikePlugin) runDocker(command string, args ...string) (output []byte, err error){
	docker, err := exec.LookPath("docker")
    if err != nil {
        log.Errorf("error finding docker binary: %s",err)
        return nil,err
    } 
    dockerArgs := []string{
    	"-H="+p.interlockConfig.SwarmUrl,
    	"--tlsverify=1",
    	"--tlscacert="+p.interlockConfig.TLSCaCert,
    	"--tlscert="+p.interlockConfig.TLSCert,
    	"--tlskey="+p.interlockConfig.TLSKey,
    	command,
    }
    dockerArgs =append(dockerArgs,args...)
    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("DOCKER ARGS  %+v",dockerArgs))
    return exec.Command(docker,dockerArgs...).Output()
    

}
func (p AerospikePlugin) runAsinfoTip(args ...string) bool{
	asinfo, err := exec.LookPath("asinfo")
    if err != nil{ 
    	log.Errorf("error finding asinfo binary: %s", err)
        return false 
    }
    time.Sleep(time.Second*5)	//sleep 5s for ASD to be ready
    cmd := exec.Command(asinfo,args...)
    plugins.Log(pluginInfo.Name, log.DebugLevel,fmt.Sprintf("asinfo COMMANMD: %+v", cmd))
    if output,err := cmd.Output(); err != nil {
        log.Errorf("error running asinfo: %s",err)
        plugins.Log(pluginInfo.Name, log.DebugLevel,fmt.Sprintf("asinfo output: %+v", output))
        return false
    }
    plugins.Log(pluginInfo.Name, log.DebugLevel,fmt.Sprintf("Sending asinfo: %s %+v", asinfo, args))
    return true   
}
func (p AerospikePlugin) clusterAerospike(eventId string) error {

    // scan all running containers
    filter := "label=com.aerospike.cluster="+p.pluginConfig.ClusterName
    args := []string{"-q","--filter",filter}
    filterOut, err := p.runDocker("ps",args...)

    if err != nil  {
        log.Errorf("error running docker: %s",err)
        return err
    }
    containers := strings.Split(string(filterOut),"\n")
    if ! p.contains(containers, eventId) {
    	plugins.Log(pluginInfo.Name, log.InfoLevel, fmt.Sprintf("CONTAINER %s not part of Aerospike Cluster",eventId))
    	return nil
    }
    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("Containers= %s", containers))


	for _,cnt := range containers {
		plugins.Log(pluginInfo.Name, log.DebugLevel,fmt.Sprintf("ContainerID=%v EventID=%v",cnt, eventId))
	

    	if cnt == eventId || cnt == ""{
            continue
        }
    	cInfo,_ := p.client.InspectContainer(cnt) // missing network info

	    inspectArgs := []string{
	    	"-f",
	    	"{{.NetworkSettings.Networks."+p.pluginConfig.NetworkName+".IPAddress}}",
	    	eventId,
	    }
	    cmdOut, err := p.runDocker("inspect",inspectArgs...)
	    if err != nil {
	        log.Errorf("error running docker: %s",err)
	        return err
	    }  
	    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("DOCKER OUTPUT %+v",string(cmdOut)))

	    // strip trailing newline
	    ip := string(cmdOut)
	    ipLen := len(ip)-1
	    if ipLen > 0 {
	    	ip = ip[:ipLen]
	    }
	    if len(cInfo.Name) > 0{
		    infoArgs := []string{
		    	"-h",
		    	cInfo.Name[1:],    //strip leading / from cInfo.Name
		    	"-v",
		    	"tip:host="+ip[:ipLen]+";port="+p.pluginConfig.MeshPort, 
		    }
			plugins.Log(pluginInfo.Name, log.InfoLevel, fmt.Sprintf("ARGS  %+v",infoArgs))	
			p.runAsinfoTip(infoArgs...)
		}else{
			log.Errorf("Container name not found")
		}
    }
	return nil
}

func (p AerospikePlugin) Init() error {
	return nil
}

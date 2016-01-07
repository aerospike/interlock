package aerospike

import (
	"fmt"
	"time"
    "os"
    "strings"
    "os/exec"
    "encoding/json"
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
/*	plugins.Log(pluginInfo.Name, log.InfoLevel,
		fmt.Sprintf("action=received event=%s time=%d",
			event.Id,
			event.Time,
		),
	)
*/
    switch event.Status {
    case "start":
        if err := p.clusterAerospike(event); err != nil{
            return err
        }
    }
    return nil
}

func contains(events []dockerclient.Container, event *dockerclient.Event) bool {
	for _,a := range events {
		if a.Id == event.Id{
			return true
		}
	}
	return false
}
/*
func findIP(b []byte, myNetwork string) (ip string ,err error) {
	var result []map[string]interface{}
    err = json.Unmarshal(b,&result); 
    networkSettings := result[0]["NetworkSettings"]
    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("JSON OUTPUT %+v",networkSettings))
    mapNetworkSettings := networkSettings.(map[string]interface{})
    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("JSON2 MAP %+v",mapNetworkSettings["Networks"]))
    mapDocker := mapNetworkSettings["Networks"].(map[string]interface{})
    if mapDocker[myNetwork] == nil{
    	err = errors.New("Network not found: "+myNetwork)
    	return
    }
    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("DOCKER MAP %+v",mapDocker[myNetwork]))
    mapDockerIP := mapDocker[myNetwork].(map[string]interface{})
    ip = mapDockerIP["IPAddress"].(string)
    plugins.Log(pluginInfo.Name, log.InfoLevel, fmt.Sprintf("DOCKER IP %+v",ip))
    return 
}
*/
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
func (p AerospikePlugin) clusterAerospike(event *dockerclient.Event) error {

    // scan all running containers
    plugins.Log(pluginInfo.Name, log.InfoLevel,fmt.Sprintf("status= %s",event.Status))

    filter := "com.aerospike.cluster="+p.pluginConfig.ClusterName
    filterMap := map[string][]string{ "label":[]string{filter} }
    mapF, _ := json.Marshal(filterMap)
    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("FILTER= %s", string(mapF)))
    containers, err := p.client.ListContainers(false,false,string(mapF)) //filter based on aerospike on 3rd param
	
    if err != nil  {
        return err
    }
    if ! contains(containers, event) {
    	plugins.Log(pluginInfo.Name, log.InfoLevel, fmt.Sprintf("CONTAINER not part of Aerospike Cluster"))
    	return nil
    }
    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("Containers= %s", containers))

 //   network,_ := p.client.InspectNetwork(p.pluginConfig.NetworkName)

    filter = "label="+filter
    args := []string{"-q","--filter",filter}
   	filterOut, err := p.runDocker("ps",args...)
   	if err != nil {
	    log.Errorf("error running docker: %s",err)
	    return err
	}  
	temp := strings.Split(string(filterOut),"\n")

//    newNode,_ := p.client.InspectContainer(event.Id)
    // for each container 
 //   for _, cnt := range containers {
//    	plugins.Log(pluginInfo.Name, log.DebugLevel,fmt.Sprintf("ContainerID=%v EventID=%v",cnt.Id, event.Id))
	for _,cnt := range temp {
		plugins.Log(pluginInfo.Name, log.DebugLevel,fmt.Sprintf("ContainerID=%v EventID=%v",cnt, event.Id[:12]))
	
 //       if cnt.Id == event.Id {
    	if cnt == event.Id[:12] || cnt == ""{
            continue
        }
//    	cInfo,_ := p.client.InspectContainer(cnt.Id[:12]) // missing network info
    	cInfo,_ := p.client.InspectContainer(cnt) // missing network info

	    inspectArgs := []string{
	    	"-f",
	    	"{{.NetworkSettings.Networks."+p.pluginConfig.NetworkName+".IPAddress}}",
	    	event.Id,
	    }
	    cmdOut, err := p.runDocker("inspect",inspectArgs...)
	    if err != nil {
	        log.Errorf("error running docker: %s",err)
	        return err
	    }  
	    plugins.Log(pluginInfo.Name, log.DebugLevel, fmt.Sprintf("DOCKER OUTPUT %+v",string(cmdOut)))
/*
	    ip2, err := findIP(cmdOut,p.pluginConfig.NetworkName)
	    if err != nil {
	        log.Errorf("error running docker: %s",err)
	        return err
	    } 
*/
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

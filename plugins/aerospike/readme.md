# Aerospike

The Aerospike plugin adds an event driven manager that would cluster aerospike containers.

This works by listening on the event stream for containers being created. When a container is created, the plugin will determine if it is part of the Aerospike Cluster.
The plugin will take care of forming the containers into a cluster.

# Configuration
The following configuration is available through environment variables:

- `AEROSPIKE_CLUSTER_NAME`: The value of the label `com.aerospike.cluster` that signifies a cluster (default:aerospike)
- `AEROSPIKE_NETWORK_NAME`: The docker overlay network on which aerospike will cluster on (default: docker)
- `AEROSPIKE_MESH_PORT`: The port that aerospike forms a mesh on (default: 3002)


# Machine
If you used Machine to create your Swarm, you can use this command to start Interlock:

    docker run -d --net prod -e AEROSPIKE_NETWORK_NAME=prod  --rm  -v /var/lib/boot2docker:/etc/docker  rguo/nterlock --swarm-url=$DOCKER_HOST --swarm-tls-ca-cert=/etc/docker/ca.pem --swarm-tls-cert=/etc/docker/server.pem --swarm-tls-key=/etc/docker/server-key.pem --debug -p aerospike start

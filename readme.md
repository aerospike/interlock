# Interlock
Dynamic, event-driven Docker plugin system using [Swarm](https://github.com/docker/swarm).

# Usage
Run `docker run ehazlett/interlock list-plugins` to show available plugins.

Example:

`docker run -P ehazlett/interlock -s tcp://1.2.3.4:2375 --plugin example start`

# Commandline options

- `--swarm-url`: url to swarm (default: tcp://127.0.0.1:2375)
- `--swarm-tls-ca-cert`: TLS CA certificate to use with swarm (optional)
- `--swarm-tls-cert`: TLS certificate to use with swarm (optional)
- `--swarm-tls-key`: TLS certificate key to use with swarm (options)
- `--plugin`: enable plugin
- `--debug`: enable debug output
- `--version`: show version and exit

# Plugins
See the [Plugins](https://github.com/ehazlett/interlock/tree/master/plugins)
directory for available plugins and their corresponding readme.md for usage.

| Name | Description |
|-----|-----|
| [Example](https://github.com/ehazlett/interlock/tree/master/plugins/example) | Example Plugin for Reference|
| [HAProxy](https://github.com/ehazlett/interlock/tree/master/plugins/haproxy) | [HAProxy](http://www.haproxy.org/) Load Balancer |
| [Nginx](https://github.com/ehazlett/interlock/tree/master/plugins/nginx) | [Nginx](http://nginx.org) Load Balancer |
| [Stats](https://github.com/ehazlett/interlock/tree/master/plugins/stats) | Container stat forwarding to [Carbon](http://graphite.wikidot.com/carbon) |
| [Aerospike](https://github.com/aerospike/interlock/tree/master/plugins/aerospike) | [Aerospike](http://aerospike.com) database cluster tipper |

# Building
To build a local copy of Interlock, you must have the following:

- Go 1.5+
- Use the Go vendor experiment

You can use the `Makefile` to build the binary.  For example:

`make build`

This will build the binary in `interlock/interlock`.

There is also a Docker image target in the makefile.  You can build it with
`make image`.

You can also use Docker to build in a container if you do not want to worry
about the host Go setup.  To build in a container run:

`make build-container`

This will build the executable and place in `interlock/interlock`.  Note: this
executable will be built for Linux so you will either need to build a container
afterword or be using Linux as your host OS to use.

# License
Licensed under the Apache License, Version 2.0. See LICENSE for full license text.

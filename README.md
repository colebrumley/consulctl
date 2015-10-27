# consulctl
Command line client for the Consul HTTP API

I started writing `consulctl` because there isn't a feature-complete command line client for Consul that supports TLS and basic auth available right now.  That's all fine and good if you like crafting long tedious curl commands, but some of us are lazy and prefer setting flags to URL queries.

This project is nowhere near complete yet, but most of the API is covered with at least list and get/set method commands where appropriate. The CLI help should be all you need, let me know if something is broken

```
NAME:
   consulctl - Command line client for Consul

USAGE:
   consulctl [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   acl		Manipulate the ACL catalog
   backup	Dump Consul's KV and Service databases to JSON
   check	Manipulate the health check catalog
   event	View or fire events
   kv, store	Manipulate the key-value store
   restore	Restore a JSON backup
   agent	Manipulate the current agent
   service	Manipulate the service catalog
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --cacert, -r 			(Optional) SSL client CA cert [$CONSULCTL_CACERT]
   --cert, -c 				(Optional) SSL client cert [$CONSULCTL_CERT]
   --key, -k 				(Optional) SSL client key [$CONSULCTL_KEY]
   --datacenter, -d "dc1"		Datacenter [$CONSULCTL_DATACENTER]
   --addr, -a "http://127.0.0.1:8500"	Consul API address [$CONSULCTL_ADDR]
   --insecure, -i			(Optional) Skip SSL host verification [$CONSULCTL_INSECURE]
   --username, -n 			(Optional) HTTP Basic auth user [$CONSULCTL_USERNAME]
   --password, -p 			(Optional) HTTP Basic auth password [$CONSULCTL_PASSWORD]
   --token, -t 				(Optional) Consul ACL Token [$CONSULCTL_TOKEN]
   --verbose, -j			Use verbose output (usually means JSON) [$CONSULCTL_VERBOSE]
   --help, -h				show help
   --version, -v			print the version
```

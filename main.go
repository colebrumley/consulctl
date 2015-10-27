package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "consulctl"
	app.Usage = "Command line client for Consul"
	app.Version = "0.0.1"
	app.HelpName = "consulctl"
	app.Flags = GlobalFlags

	app.Commands = []cli.Command{
		ACLCommand,
		BackupCommand,
		CheckCommand,
		EventsCommand,
		KvCommand,
		RestoreCommand,
		AgentCommand,
		ServiceCommand,
	}

	app.Run(os.Args)
}

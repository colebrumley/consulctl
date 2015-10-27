package main

import (
	"github.com/codegangsta/cli"
)

var (
	GlobalFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "cacert,r",
			Usage:  "(Optional) SSL client CA cert",
			EnvVar: "CONSULCTL_CACERT",
		},
		cli.StringFlag{
			Name:   "cert,c",
			Usage:  "(Optional) SSL client cert",
			EnvVar: "CONSULCTL_CERT",
		},
		cli.StringFlag{
			Name:   "key,k",
			Usage:  "(Optional) SSL client key",
			EnvVar: "CONSULCTL_KEY",
		},
		cli.StringFlag{
			Name:   "datacenter,d",
			Usage:  "Datacenter",
			Value:  "dc1",
			EnvVar: "CONSULCTL_DATACENTER",
		},
		cli.StringFlag{
			Name:   "addr,a",
			Value:  "http://127.0.0.1:8500",
			Usage:  "Consul API address",
			EnvVar: "CONSULCTL_ADDR",
		},
		cli.BoolFlag{
			Name:   "insecure,i",
			Usage:  "(Optional) Skip SSL host verification",
			EnvVar: "CONSULCTL_INSECURE",
		},
		cli.StringFlag{
			Name:   "username,n",
			Usage:  "(Optional) HTTP Basic auth user",
			EnvVar: "CONSULCTL_USERNAME",
		},
		cli.StringFlag{
			Name:   "password,p",
			Usage:  "(Optional) HTTP Basic auth password",
			EnvVar: "CONSULCTL_PASSWORD",
		},
		cli.StringFlag{
			Name:   "token,t",
			Usage:  "(Optional) Consul ACL Token",
			EnvVar: "CONSULCTL_TOKEN",
		},
		cli.BoolFlag{
			Name:   "verbose,j",
			Usage:  "Use verbose output (usually means JSON)",
			EnvVar: "CONSULCTL_VERBOSE",
		},
	}

	BackupCommand = cli.Command{
		Name:      "backup",
		Usage:     "Dump Consul's KV and Service databases to JSON",
		ArgsUsage: " ",
		Action:    BackupConsul,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "outfile,o",
				Usage: "Write output to a file",
			},
			cli.BoolFlag{
				Name:  "indent,i",
				Usage: "Create indented JSON",
			},
		},
	}

	RestoreCommand = cli.Command{
		Name:      "restore",
		Usage:     "Restore a JSON backup",
		ArgsUsage: "restore_file.json",
		Action:    RestoreConsul,
	}

	AgentCommand = cli.Command{
		Name:      "agent",
		Usage:     "Manipulate the current agent",
		ArgsUsage: " ",
		Subcommands: []cli.Command{
			{
				Name:      "members",
				Usage:     "Get info about cluster members",
				ArgsUsage: " ",
				Action:    GetMemberInfo,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "wan,w",
						Usage: "Include WAN members",
					},
				},
			},
			cli.Command{
				Name:      "info",
				Usage:     "Get info about the current agent",
				ArgsUsage: " ",
				Action:    GetSelfInfo,
			},
			cli.Command{
				Name:      "maintenance",
				Aliases:   []string{"maint"},
				Usage:     "Toggle maintenance mode on the current node",
				ArgsUsage: "[enable|disable]",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "reason,r",
						Usage: "Set an optional reason for maintenance mode",
					},
				},
				Action: ToggleMaintenanceMode,
			},
		},
	}

	EventsCommand = cli.Command{
		Name:      "event",
		Usage:     "View or fire events",
		ArgsUsage: " ",
		Subcommands: []cli.Command{
			cli.Command{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List recent events",
				ArgsUsage: "optional-name",
				Action:    ListEvents,
			},
			cli.Command{
				Name:      "fire",
				Aliases:   []string{"new"},
				Usage:     "Fire a new event",
				ArgsUsage: "name payload",
				Action:    FireEvent,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "node,n",
						Usage: "Node filter regex",
					},
					cli.StringFlag{
						Name:  "service,s",
						Usage: "Service filter regex",
					},
					cli.StringFlag{
						Name:  "tag,t",
						Usage: "Tag filter regex",
					},
				},
			},
		},
	}

	KvCommand = cli.Command{
		Name:    "kv",
		Aliases: []string{"store"},
		Usage:   "Manipulate the key-value store",
		Subcommands: []cli.Command{
			{
				Name:      "get",
				Usage:     "Get key(s)",
				ArgsUsage: "/my/key1 /my/key2...",
				Action:    GetKvKey,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "recurse,r",
						Usage: "Get keys recursively",
					},
				},
			},
			{
				Name:      "set",
				Usage:     "Set a key's value",
				ArgsUsage: "/my/key 'my value'",
				Action:    SetKvKey,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "quiet,q",
						Usage: "Suppress confirmation message",
					},
				},
			},
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List keys",
				ArgsUsage: "/optional/root",
				Action:    ListKeys,
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  "recurse,r",
						Usage: "Get keys recursively",
					},
				},
			},
		},
	}

	ServiceCommand = cli.Command{
		Name:      "service",
		Usage:     "Manipulate the service catalog",
		ArgsUsage: " ",
		Subcommands: []cli.Command{
			{
				Name:      "register",
				Aliases:   []string{"new"},
				Usage:     "Create or edit a service",
				ArgsUsage: "name=$ address=$ port=$ tags=$,$,$",
				Action:    SetService,
			},
			{
				Name:      "deregister",
				Aliases:   []string{"rm"},
				Usage:     "Remove a service",
				ArgsUsage: "service-id",
				Action:    DeleteService,
			},
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List service info",
				ArgsUsage: " ",
				Action:    ListServices,
			},
			{
				Name:      "maintenance",
				Aliases:   []string{"maint"},
				Usage:     "Toggle maintenance mode",
				ArgsUsage: "[enable|disable] service 'Optional reason string'",
				Action:    ServiceMaint,
			},
		},
	}

	ACLCommand = cli.Command{
		Name:      "acl",
		Usage:     "Manipulate the ACL catalog",
		ArgsUsage: " ",
		Subcommands: []cli.Command{
			{
				Name:      "register",
				Aliases:   []string{"new"},
				Usage:     "Create an ACL token",
				ArgsUsage: "'{\"rule\":\"json\"}'",
				Action:    CreateACL,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "type,t",
						Usage: "ACL Token type (client or management)",
					},
					cli.StringFlag{
						Name:  "name,n",
						Usage: "ACL name",
					},
					cli.StringFlag{
						Name:  "id,i",
						Usage: "ID of ACL",
					},
				},
			},
			{
				Name:      "deregister",
				Aliases:   []string{"rm"},
				Usage:     "Remove an ACL token",
				ArgsUsage: "service-id",
				Action:    DeregisterACL,
			},
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List ACL tokens",
				ArgsUsage: " ",
				Action:    ListACL,
			},
			{
				Name:      "clone",
				Usage:     "Clone a token into a new entry",
				ArgsUsage: " ",
				Action:    CloneACL,
			},
			{
				Name:      "update",
				Aliases:   []string{"up"},
				Usage:     "Update an ACL token",
				ArgsUsage: "'{\"rule\":\"json\"}'",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "type,t",
						Usage: "ACL Token type (client or management)",
					},
					cli.StringFlag{
						Name:  "name,n",
						Usage: "ACL name",
					},
					cli.StringFlag{
						Name:  "id,i",
						Usage: "ID of ACL (required)",
					},
				},
				Action: UpdateACL,
			},
		},
	}

	CheckCommand = cli.Command{
		Name:      "check",
		Usage:     "Manipulate the health check catalog",
		ArgsUsage: " ",
		Subcommands: []cli.Command{
			{
				Name:      "register",
				Aliases:   []string{"new"},
				Usage:     "Create or edit a check",
				ArgsUsage: "checkname",
				Action:    RegisterCheck,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "http,h",
						Usage: "HTTP check endpoint",
					},
					cli.StringFlag{
						Name:  "tcp,t",
						Usage: "TCP check endpoint",
					},
					cli.StringFlag{
						Name:  "interval,i",
						Usage: "Check interval",
					},
					cli.StringFlag{
						Name:  "script,x",
						Usage: "Check script path",
					},
					cli.StringFlag{
						Name:  "ttl,l",
						Usage: "Check TTL",
					},
					cli.StringFlag{
						Name:  "service,s",
						Usage: "Service ID to associate with this check",
					},
					cli.StringFlag{
						Name:  "notes,n",
						Usage: "Notes about this check",
					},
					cli.StringFlag{
						Name:  "timeout,o",
						Usage: "HTTP/TCP timeout",
					},
				},
			},
			{
				Name:      "deregister",
				Aliases:   []string{"rm"},
				Usage:     "Remove a check",
				ArgsUsage: "service-id",
				Action:    DeregisterCheck,
			},
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "List health check info",
				ArgsUsage: " ",
				Action:    ListChecks,
			},
			{
				Name:    "update",
				Aliases: []string{"up"},
				Usage:   "Update TTL check",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "note,n",
						Usage: "Notes for this update",
					},
				},
				ArgsUsage: "check-id [pass|fail|warn]",
				Action:    UpdateCheck,
			},
		},
	}
)

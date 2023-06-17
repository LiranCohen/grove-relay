package flags

import "github.com/urfave/cli/v2"

var Host = cli.StringFlag{
	Name:        "host",
	Usage:       "service host",
	DefaultText: "0.0.0.0",
}

var Port = cli.IntFlag{
	Name:        "port",
	Usage:       "service port",
	DefaultText: "7447",
}

var Postgres = cli.StringFlag{
	Name:     "pgconn",
	Usage:    "postgres connection string",
	Required: true,
}

var Name = cli.StringFlag{
	Name:     "name",
	Usage:    "name",
	Required: true,
}

var Description = cli.StringFlag{
	Name:  "description",
	Usage: "description, shows up o nip11",
}

var PubKey = cli.StringFlag{
	Name:  "pubkey",
	Usage: "owner's public key, shows up on nip11",
}

var Contact = cli.StringFlag{
	Name:  "contact",
	Usage: "owner's contact email, shows up on nip11",
}

var MaxEvent = cli.IntFlag{
	Name:        "maxeventsz",
	Usage:       "max event size in bytes",
	DefaultText: "100000",
}

var MaxCache = cli.IntFlag{
	Name:        "maxcachesz",
	Usage:       "max cache entries, 10000 max",
	DefaultText: "100",
}

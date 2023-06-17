package flags

import "github.com/urfave/cli/v2"

var Host = cli.StringFlag{
	Name:        "host",
	Usage:       "service host [default 0.0.0.0]",
	DefaultText: "0.0.0.0",
}

var Port = cli.IntFlag{
	Name:        "port",
	Usage:       "service port [default 7447]",
	DefaultText: "7447",
}

var Postgres = cli.StringFlag{
	Name:     "pgconn",
	Usage:    "Postgres Connection String",
	Required: true,
}

var Name = cli.StringFlag{
	Name:     "name",
	Usage:    "Name",
	Required: true,
}

var Description = cli.StringFlag{
	Name:  "description",
	Usage: "Description",
}

var PubKey = cli.StringFlag{
	Name:  "pubkey",
	Usage: "Public Key",
}

var Contact = cli.StringFlag{
	Name:  "contact",
	Usage: "Contact",
}

var MaxEvent = cli.IntFlag{
	Name:  "maxeventsz",
	Usage: "Max Event Size",
}

var MaxCache = cli.IntFlag{
	Name:  "maxcachesz",
	Usage: "Max Cache Size",
}

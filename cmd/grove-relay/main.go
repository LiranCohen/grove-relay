package main

import (
	"log"
	"os"

	"github.com/Bitcoin-Grove/grove-relay/cmd/flags"
	"github.com/Bitcoin-Grove/grove-relay/server"
	"github.com/fiatjaf/relayer/v2"
	"github.com/fiatjaf/relayer/v2/storage/postgresql"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "relay",
		Usage: "a whitelisted relay for bitcoin grove",
		Flags: []cli.Flag{
			&flags.Host,
			&flags.Port,
			&flags.Postgres,
			&flags.Name,
			&flags.Description,
			&flags.Software,
			&flags.PubKey,
			&flags.Contact,
			&flags.MaxEvent,
			&flags.MaxCache,
		},
		Action: func(c *cli.Context) error {
			storage := &postgresql.PostgresBackend{DatabaseURL: c.String("pgconn")}

			s := server.New(
				server.WithStorage(storage),
				server.WithName(c.String("name")),
				server.WithDescription(c.String("description")),
				server.WithPubKey(c.String("pubkey")),
				server.WithContact(c.String("contact")),
				server.WithSoftware("software"),
				server.WithMaxEventSize(100000),
				server.WithMaxCache(1000),
			)

			if server, err := relayer.NewServer(s); err != nil {
				return err
			} else {
				return server.Start(c.String("host"), c.Int("port"))
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

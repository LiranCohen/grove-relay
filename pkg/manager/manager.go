package manager

import (
	"fmt"
	"log"
	"time"

	"github.com/fiatjaf/relayer/v2/storage/postgresql"
)

func New(storage *postgresql.PostgresBackend, internalAdmins ...string) (*Admin, error) {
	a := &Admin{
		internalAdmins,
		storage,
	}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

type Admin struct {
	internalAdmins []string

	storage *postgresql.PostgresBackend
}

func (a *Admin) init() error {
	stmt := `CREATE TABLE IF NOT EXISTS admin (
    pubkey CHAR(64) NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (pubkey)
	);`
	_, err := a.storage.Exec(stmt)
	return err
}

func (a *Admin) IsAdmin(pubKey string) bool {
	stmt := `SELECT COUNT(*) 
	FROM admin
	WHERE pubkey = $1 AND (deleted IS NULL OR deleted < $2)
	`
	var count int
	err := a.storage.QueryRow(stmt, pubKey, time.Now()).Scan(&count)
	if err != nil {
		log.Printf("[DB ERROR] could not query for pubkey %s - %s\n", pubKey, err.Error())
		return false
	}

	return count > 0
}

func (a *Admin) Add(pubKey string) error {
	stmt := `
	INSERT INTO admin (pubkey, deleted) 
	VALUES ($1, NULL)
	ON CONFLICT (pubkey) DO UPDATE 
	SET deleted = NULL
	`

	if _, err := a.storage.Exec(stmt, pubKey); err != nil {
		return fmt.Errorf("could not insert pubkey: %w", err)
	}

	return nil
}

func (a *Admin) Remove(pubKey string) error {
	stmt := `
	UPDATE admin
	SET deleted = CURRENT_TIMESTAMP
	WHERE pubkey = $1
	`

	if _, err := a.storage.Exec(stmt, pubKey); err != nil {
		return fmt.Errorf("could not remove pubkey: %w", err)
	}

	return nil
}

package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fiatjaf/relayer/v2"
	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/nbd-wtf/go-nostr"
)

func New(storage *postgresql.PostgresBackend, internalAdmins ...string) (*Admin, error) {

	a := &Admin{
		storage:        storage,
		internalAdmins: map[string]struct{}{},
	}

	for _, pubKey := range internalAdmins {
		a.internalAdmins[pubKey] = struct{}{}
	}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

type Admin struct {
	internalAdmins map[string]struct{}

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
	if _, ok := a.internalAdmins[pubKey]; ok {
		return true
	}

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

func (a *Admin) add(pubKeys ...string) error {
	if len(pubKeys) == 0 {
		return nil
	}

	stmtValues := []string{}
	values := []any{}
	for n, pubKey := range pubKeys {
		stmtValues = append(stmtValues, fmt.Sprintf("(%d, NULL)", n+1))
		values = append(values, pubKey)
	}

	stmt := "INSERT INTO admin (pubkey, deleted) VALUES " + strings.Join(stmtValues, ", ") + "ON CONFLICT (pubkey) DO UPDATE SET deleted = NULL;"
	if _, err := a.storage.Exec(stmt, values...); err != nil {
		return fmt.Errorf("could not insert pubkey: %w", err)
	}

	return nil
}

func (a *Admin) remove(pubKeys ...string) error {
	if len(pubKeys) == 0 {
		return nil
	}

	stmtValues := []string{}
	values := []any{}
	for n, pubKey := range pubKeys {
		stmtValues = append(stmtValues, fmt.Sprintf("pubkey = %d", n+1))
		values = append(values, pubKey)
	}
	stmt := "UPDATE admin SET deleted = CURRENT_TIMESTAMP WHERE " + strings.Join(stmtValues, "OR ") + ";"
	if _, err := a.storage.Exec(stmt, values...); err != nil {
		return fmt.Errorf("could not remove pubkey: %w", err)
	}

	return nil
}

func (a *Admin) HandleAdminType(ws *relayer.WebSocket, request []json.RawMessage) {
	var notice string
	defer func() {
		if notice != "" {
			ws.WriteJSON(nostr.NoticeEnvelope(notice))
		}
	}()

	var evt nostr.Event
	if err := json.Unmarshal(request[1], &evt); err != nil {
		notice = "failed to decode auth event: " + err.Error()
		return
	}

	if ok, err := evt.CheckSignature(); !ok || err != nil {
		notice = "invalid signature"
		if err != nil {
			notice += ": " + err.Error()
		}
		return
	}

	if !a.IsAdmin(evt.PubKey) {
		notice = "failed to auth"
		return
	}

	toAdd := evt.Tags.GetAll([]string{"add"})
	fmt.Printf("To Add: %+v\n", toAdd)
	toRemove := evt.Tags.GetAll([]string{"remove"})
	fmt.Printf("To Add: %+v\n", toRemove)

}

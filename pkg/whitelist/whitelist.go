package whitelist

import (
	"container/list"
	"fmt"
	"log"
	"time"

	"github.com/fiatjaf/relayer/v2/storage/postgresql"
)

type Cache struct {
	maxCapacity int
	storage     *postgresql.PostgresBackend

	list  *list.List
	items map[interface{}]*list.Element
}

func WithMaxCapacity(max int) Option {
	return func(w *Cache) error {
		if max > 10000 {
			log.Printf("[Max Capacity Exceeded] %d too high, resetting to max 10000\n", max)
			w.maxCapacity = 10000
		} else {
			w.maxCapacity = max
		}

		return nil
	}
}

func WithStorage(storage *postgresql.PostgresBackend) Option {
	return func(w *Cache) error {
		w.storage = storage
		return nil
	}
}

type Option func(*Cache) error

func New(opts ...Option) *Cache {
	w := &Cache{}

	for _, opt := range opts {
		if err := opt(w); err != nil {
			log.Fatal(fmt.Errorf("could not load option: %w", err))
		}
	}

	if w.storage == nil {
		log.Fatalln("[Invalid Whitelist Storage]")
	}

	if w.maxCapacity < 10 {
		log.Printf("[Max Capacity Min] set to 10")
	}

	if err := w.init(); err != nil {
		log.Fatal(fmt.Errorf("could not init: %w", err))
	}

	return w
}

func (w *Cache) init() error {
	stmt := `CREATE TABLE IF NOT EXIST relay_whitelist (
    pubkey CHAR(64) NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted TIMESTAMP WITH TIME ZONE,
    PRIMARY KEY (pubkey)
	);`
	_, err := w.storage.Exec(stmt)
	return err
}

func (w *Cache) Allowed(pubKey string) bool {
	if el, ok := w.items[pubKey]; ok {
		w.list.MoveToFront(el)
		return true
	}

	stmt := `
		SELECT COUNT(*) 
		FROM relay_whitelist 
		WHERE pubkey = $1 AND (deleted IS NULL OR deleted < $2)
	`
	var count int
	err := w.storage.QueryRow(stmt, pubKey, time.Now()).Scan(&count)
	if err != nil {
		log.Printf("[DB ERROR] could not query for pubkey %s - %s\n", pubKey, err.Error())
		return false
	}
	if count > 0 {
		w.setCache(pubKey)
		return true
	} else {
		return false
	}
}

func (w *Cache) SetAllowed(pubKey string) error {
	stmt := `
		INSERT INTO relay_whitelist (pubkey, deleted) 
		VALUES ($1, NULL)
		ON CONFLICT (pubkey) DO UPDATE 
		SET deleted = NULL
	`

	if _, err := w.storage.Exec(stmt, pubKey); err != nil {
		return fmt.Errorf("could not insert pubkey: %w", err)
	}

	w.setCache(pubKey)
	return nil
}

func (w *Cache) Deactivate(pubKey string) error {
	stmt := `
		UPDATE relay_whitelist
		SET deleted = CURRENT_TIMESTAMP
		WHERE pubkey = $1
	`
	if _, err := w.storage.Exec(stmt, pubKey); err != nil {
		return fmt.Errorf("could not remove pubkey: %w", err)
	}

	w.delCache(pubKey)
	return nil
}

func (w *Cache) setCache(pubKey string) {
	w.cacheClean()
	if el, ok := w.items[pubKey]; ok {
		w.list.MoveToFront(el)
		return
	}
	el := w.list.PushFront(pubKey)
	w.items[pubKey] = el
}

func (w *Cache) cacheClean() {
	if w.list.Len() >= w.maxCapacity {
		el := w.list.Back()
		w.list.Remove(el)
		delete(w.items, el)
	}
}

func (w *Cache) delCache(pubKey string) {
	delete(w.items, pubKey)
	w.cacheClean()
}

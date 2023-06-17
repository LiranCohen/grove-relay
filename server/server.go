package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Bitcoin-Grove/grove-relay/pkg/whitelist"

	"github.com/fiatjaf/relayer/v2"
	"github.com/fiatjaf/relayer/v2/storage/postgresql"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip11"
)

type Option func(*Relay) error

func WithName(name string) Option {
	return func(r *Relay) error {
		r.name = name
		return nil
	}
}

func WithDescription(description string) Option {
	return func(r *Relay) error {
		r.description = description
		return nil
	}
}

func WithPubKey(pubKey string) Option {
	return func(r *Relay) error {
		r.pubKey = pubKey
		return nil
	}
}

func WithContact(contact string) Option {
	return func(r *Relay) error {
		r.contact = contact
		return nil
	}
}

func WithSoftwareURL(software string) Option {
	return func(r *Relay) error {
		r.software = software
		return nil
	}
}

func WithMaxEventSize(max int) Option {
	return func(r *Relay) error {
		r.maxEventSize = max
		return nil
	}
}

func WithMaxCache(max int) Option {
	return func(r *Relay) error {
		r.maxCacheSize = max
		return nil
	}
}

func WithStorage(storage *postgresql.PostgresBackend) Option {
	return func(r *Relay) error {
		r.storage = storage
		return nil
	}
}

type Relay struct {
	name          string
	description   string
	pubKey        string
	contact       string
	software      string
	maxEventSize  int
	maxCacheSize  int
	supportedNIPs []int

	whitelist *whitelist.Cache
	storage   *postgresql.PostgresBackend
}

func New(opts ...Option) *Relay {
	r := &Relay{}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			log.Default().Panic(err)
		}
	}

	if r.storage == nil {
		log.Default().Panic(fmt.Errorf("could not load storage - empty interface"))
	}

	if r.name == "" {
		r.name = "WhitelistRelay"
	}

	if r.description == "" {
		r.description = "A whitelist relay based on fiatjaf's relayer."
	}

	if r.software == "" {
		r.software = "https://github.com/fiatjaf/relayer"
	}

	r.supportedNIPs = uniqueInts(r.supportedNIPs, []int{9, 11, 12, 15, 16, 20, 33, 42})

	return r
}

// ChatGPT Function
func uniqueInts(list1, list2 []int) []int {
	uniqueMap := make(map[int]struct{})
	for _, num := range list1 {
		uniqueMap[num] = struct{}{}
	}
	for _, num := range list2 {
		uniqueMap[num] = struct{}{}
	}

	var uniqueSlice []int
	for num := range uniqueMap {
		uniqueSlice = append(uniqueSlice, num)
	}
	return uniqueSlice
}

func (r *Relay) Name() string {
	return r.name
}

func (r *Relay) Storage(ctx context.Context) relayer.Storage {
	return r.storage
}

func (r *Relay) Init() error {
	if r.whitelist == nil {
		r.whitelist = whitelist.New(
			whitelist.WithMaxCapacity(1000),
			whitelist.WithStorage(r.storage),
		)
	}
	return nil
}

func (r *Relay) AcceptEvent(ctx context.Context, evt *nostr.Event) bool {
	if !r.whitelist.Allowed(evt.PubKey) {
		return false
	}

	// block events that are too large
	jsonb, _ := json.Marshal(evt)
	return len(jsonb) <= r.maxEventSize
}

func (r *Relay) GetNIP11InformationDocument() nip11.RelayInformationDocument {
	return nip11.RelayInformationDocument{
		Name:          r.name,
		Description:   r.description,
		PubKey:        r.pubKey,
		Contact:       r.contact,
		SupportedNIPs: r.supportedNIPs,
		Software:      r.software,
		Version:       "~",
	}
}

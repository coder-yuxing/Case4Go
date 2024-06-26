package factory

import (
	"fmt"
	"log"
	"sync"
	"yuxing.code/bookstore/store"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]store.Store)
)

func Register(name string, p store.Store) {
	log.Printf("[factory]: register provider %q", name)
	providersMu.Lock()
	defer providersMu.Unlock()
	if p == nil {
		panic("store: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("store: Register called twice for provider " + name)
	}
	providers[name] = p
}

func New(providerName string) (store.Store, error) {
	providersMu.RLock()
	p, ok := providers[providerName]
	providersMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("store: unknow provider %s", providerName)
	}
	return p, nil
}

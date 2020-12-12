package memce

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tutamuniz/yasrp/minihttp/cache"
)

type MemEntryMap map[string]*cache.CacheEntry
type MemCE struct {
	sync.RWMutex
	entry MemEntryMap
}

func NewMemCE() (*MemCE, error) {

	return &MemCE{
		entry: make(MemEntryMap),
	}, nil
}

func (m *MemCE) InCache(key string) bool {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()

	if _, ok := m.entry[key]; ok {
		return true
	}
	return false
}

func (m *MemCE) Get(key string) (*cache.CacheEntry, error) {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()

	if e, ok := m.entry[key]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("Entry not found")
}

func (m *MemCE) Put(key string, item *cache.CacheEntry) error {
	m.RWMutex.Lock()
	defer m.RWMutex.Unlock()
	m.entry[key] = item
	return nil // Maybe used in the future
}

func (m *MemCE) StartEngine() {
	log.Printf("Starting Memory Cache Engine\n")
	for t := range time.Tick(time.Minute * 1) {
		log.Printf("MemCE: Running janitor %s\n", t.String())
		log.Printf("MemCE: %d itens chached \n", len(m.entry))
		m.RWMutex.Lock()
		for k, v := range m.entry {
			if v.ExpireOn.Before(time.Now()) {
				log.Printf("MemCE: Expiring %s cache\n", k)
				delete(m.entry, k)
			}
		}
		m.RWMutex.Unlock()
	}
}

/*
	InCache(string) bool
	Get(string) ([]byte, error)
	Put(string, []byte) error
	PutTTL(string, []byte) error
	StartEngine()
}
*/

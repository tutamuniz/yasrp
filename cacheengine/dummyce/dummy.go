package dummyce

import (
	"fmt"
	"log"

	"github.com/tutamuniz/yasrp/minihttp/cache"
)

type DummyCE struct{}

func NewDummyCE() (*DummyCE, error) {
	return &DummyCE{}, nil
}

func (_ DummyCE) InCache(_ string) bool {
	return false
}

func (_ DummyCE) Get(_ string) (*cache.CacheEntry, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (_ DummyCE) Put(string, *cache.CacheEntry) error {
	return fmt.Errorf("Not implemented")

}

func (_ DummyCE) StartEngine() {
	log.Printf("Starting Dummy Cache Engine\n")
}

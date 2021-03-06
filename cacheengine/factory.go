package cacheengine

import (
	"fmt"

	"github.com/tutamuniz/yasrp/cacheengine/dummyce"
	"github.com/tutamuniz/yasrp/cacheengine/memce"
	"github.com/tutamuniz/yasrp/minihttp/cache"
)

func NewCacheEngine(name string) (*cache.Cache, error) {
	var engine cache.Cache
	var err error
	switch name {
	case "dummy":
		engine, err = dummyce.NewDummyCE()
	case "mem":
		engine, err = memce.NewMemCE()
	default:
		engine = nil
		err = fmt.Errorf("Invalid Cache Engine")
	}
	return &engine, err
}

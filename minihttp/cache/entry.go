package cache

import (
	"time"

	"github.com/tutamuniz/yasrp/minihttp"
)

type CacheEntry struct {
	Req      *minihttp.Request
	Resp     *minihttp.Response
	ExpireOn time.Time
}

func MakeCacheEntry(rq *minihttp.Request, rp *minihttp.Response) *CacheEntry {
	expire := time.Now().Add(CacheMaxAge)
	return &CacheEntry{
		Req:      rq,
		Resp:     rp,
		ExpireOn: expire,
	}
}

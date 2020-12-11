package cache

type Cache interface {
	InCache(string) bool
	Get(string) ([]byte, error)
	Put(string, []byte) error
	PutTTL(string, []byte) error
	StartEngine()
}

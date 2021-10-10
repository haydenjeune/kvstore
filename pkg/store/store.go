package store

type KvStore interface {
	Get(key string) (string, bool, error)
	Set(key string, value string) error
}

package store

type InMemHashMapKVStorage struct {
	hashmap map[string]string
}

func NewInMemHashMapKVStorage() (*InMemHashMapKVStorage, error) {
	return &InMemHashMapKVStorage{
		hashmap: make(map[string]string),
	}, nil
}

func (s *InMemHashMapKVStorage) Get(key string) (string, bool, error) {
	value, exists := s.hashmap[key]
	return value, exists, nil
}

func (s *InMemHashMapKVStorage) Set(key string, value string) error {
	s.hashmap[key] = value
	return nil
}

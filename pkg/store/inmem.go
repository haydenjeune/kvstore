package store

import "github.com/haydenjeune/kvstore/pkg/bst"

type InMemSortedKVStorage struct {
	memtable *bst.BinarySearchTree
}

func NewInMemSortedKVStorage() (*InMemSortedKVStorage, error) {
	return &InMemSortedKVStorage{
		memtable: &bst.BinarySearchTree{},
	}, nil
}

func (s *InMemSortedKVStorage) Get(key string) (string, bool, error) {
	value, exists := s.memtable.Search(key)
	return value, exists, nil
}

func (s *InMemSortedKVStorage) Set(key string, value string) error {
	s.memtable.Insert(key, value)
	return nil
}

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
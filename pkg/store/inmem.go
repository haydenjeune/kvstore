package store

import "github.com/haydenjeune/kvstore/pkg/bst"

type InMemSortedKVStorage struct {
	memtable *bst.BinarySearchTree
}

func NewInMemSortedKVStorage() *InMemSortedKVStorage {
	return &InMemSortedKVStorage{
		memtable: &bst.BinarySearchTree{},
	}
}

func (s *InMemSortedKVStorage) Get(key string) (string, bool, error) {
	value, exists := s.memtable.Search(key)
	return value, exists, nil
}

func (s *InMemSortedKVStorage) Set(key string, value string) error {
	s.memtable.Insert(key, value)
	return nil
}

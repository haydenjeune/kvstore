package store

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/haydenjeune/kvstore/pkg/bst"
)

type SortedFile struct {
	index    []KeyOffset
	filename string
}

type KeyOffset struct {
	Key    string
	Offset int64
}

const RECORDS_PER_INDEX_ENTRY uint = 10
const MAX_RECORDS_PER_FILE uint = 100

func NewSortedFile(filename string) (*SortedFile, error) {
	index, err := newSparseIndexFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to build index from file '%s': %v", filename, err)
	}
	return &SortedFile{
		index:    index,
		filename: filename,
	}, nil
}

func newSparseIndexFromFile(filename string) ([]KeyOffset, error) {
	index := make([]KeyOffset, 0)

	f, err := os.Open(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return index, nil
	} else if err != nil {
		return nil, fmt.Errorf("couldn't open data file: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var offset int64 = 0
	var lastKey string
	for scanner.Scan() {
		key := strings.SplitN(scanner.Text(), ", ", 2)[0]
		if key <= lastKey && key != "" {
			return nil, fmt.Errorf("encountered out of order keys '%s' and '%s'", lastKey, key)
		}
		index = append(index, KeyOffset{Key: key, Offset: offset})
		offset += int64(len(scanner.Bytes())) + 1 // add an extra byte for the newline
		lastKey = key
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("couldn't scan data file: %v", err)
	}

	return index, nil
}

type SortedFileKvStorage struct {
	dir      string
	memtable bst.BinarySearchTree
}

func NewSortedFileKvStorage() (*SortedFileKvStorage, error) {
	return &SortedFileKvStorage{}, nil
}

func (s *SortedFileKvStorage) Get(key string) (string, bool, error) {
	value, exists := s.memtable.Search(key)
	if exists {
		return value, true, nil
	}

	return "", false, nil
}

func (s *SortedFileKvStorage) Set(key string, value string) error {
	return nil
}

package sortedfile

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/haydenjeune/kvstore/pkg/bst"
	"github.com/spf13/afero"
)

type SortedFile struct {
	index    []KeyOffset
	filename string
	fs       afero.Fs
}

type KeyOffset struct {
	Key    string
	Offset int64
}

const RECORDS_PER_INDEX_ENTRY uint = 10
const MAX_RECORDS_PER_FILE uint = 100

func NewSortedFile(filename string, fs afero.Fs) (*SortedFile, error) {
	index, err := newSparseIndexFromFile(filename, fs)
	if err != nil {
		return nil, fmt.Errorf("failed to build index from file '%s': %v", filename, err)
	}
	return &SortedFile{
		index:    index,
		filename: filename,
		fs:       fs,
	}, nil
}

func (s *SortedFile) Get(key string) (string, bool, error) {
	f, err := s.fs.Open(s.filename)
	if err != nil {
		return "", false, fmt.Errorf("failed to open sorted file '%s': %v", s.filename, err)
	}

	info, err := f.Stat()
	if err != nil {
		return "", false, fmt.Errorf("failed to get stats from file '%s': %v", s.filename, err)
	}

	l, r := getInterval(s.index, key)

	var offset int64
	if l == nil {
		offset = 0
	} else {
		offset = l.Offset
	}

	var endOffset int64
	if r != nil {
		endOffset = r.Offset
	} else {
		endOffset = info.Size()
	}

	f.Seek(offset, io.SeekStart)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() && offset < endOffset {
		line := scanner.Text()
		if strings.HasPrefix(line, key+", ") {
			value := strings.SplitN(line, ", ", 2)[1]
			return value, true, nil
		}
		offset += int64(len(scanner.Bytes()) + 1) // add one for the newline
	}
	if scanner.Err() != nil {
		return "", false, fmt.Errorf("failed to scan line at offset %d in file '%s': %v", offset, s.filename, scanner.Err())
	}
	return "", false, nil
}

func newSparseIndexFromFile(filename string, fs afero.Fs) ([]KeyOffset, error) {
	index := make([]KeyOffset, 0)

	if exists, _ := afero.Exists(fs, filename); !exists {
		return nil, fmt.Errorf("file '%s' does not exist", filename)
	}
	f, err := fs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file '%s': %v", filename, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var offset int64 = 0
	var lastKey string
	for i := uint(0); scanner.Scan(); i++ {
		key := strings.SplitN(scanner.Text(), ", ", 2)[0]
		if key <= lastKey && key != "" {
			return nil, fmt.Errorf("encountered out of order keys '%s' and '%s'", lastKey, key)
		}
		if i%RECORDS_PER_INDEX_ENTRY == 0 {
			index = append(index, KeyOffset{Key: key, Offset: offset})
		}
		offset += int64(len(scanner.Bytes())) + 1 // add an extra byte for the newline
		lastKey = key
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("couldn't scan data file: %v", err)
	}

	return index, nil
}

func writeBstToSortedFile(t *bst.BinarySearchTree, filename string, fs afero.Fs) error {
	f, err := fs.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	iter := bst.NewInOrderTraversalIterator(t)

	for iter.Next() {
		node := iter.Value()
		_, err := f.WriteString(fmt.Sprintf("%s, %s\n", node.Key(), node.Value()))
		if err != nil {
			return fmt.Errorf("failed to write line to file: %v", err)
		}
	}

	return nil
}

func getInterval(index []KeyOffset, key string) (*KeyOffset, *KeyOffset) {
	if len(index) == 0 {
		return nil, nil
	}

	if len(index) == 1 {
		if key < index[0].Key {
			return nil, &index[0]
		} else {
			return &index[0], nil
		}
	}

	mid := len(index) / 2

	var l, r *KeyOffset
	if key < index[mid].Key {
		l, r = getInterval(index[:mid], key)
		if r == nil {
			r = &index[mid]
		}
	} else {
		l, r = getInterval(index[mid:], key)
		if l == nil && mid > 0 {
			l = &index[mid-1]
		}
	}

	return l, r
}

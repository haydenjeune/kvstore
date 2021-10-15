package store

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

type HashIndexedFsAppendOnlyStorage struct {
	filename  string
	index     map[string]int64
	endOffset int64
}

func newHashIndexFromFile(filename string) (map[string]int64, error) {
	index := make(map[string]int64)

	f, err := os.Open(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return index, nil
	} else if err != nil {
		return nil, fmt.Errorf("couldn't open data file: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var offset int64 = 0
	for scanner.Scan() {
		key := strings.SplitN(scanner.Text(), ", ", 2)[0]
		index[key] = offset
		offset += int64(len(scanner.Bytes())) + 1 // add an extra byte for the newline
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("couldn't scan data file: %v", err)
	}
	return index, nil
}

func NewHashIndexedFsAppendOnlyStorage(filename string) (*HashIndexedFsAppendOnlyStorage, error) {
	// Ensure data file exists and is openable
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("couldn't open data file: %v", err)
	}

	// Get size of file to set the offset
	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("couldn't get stats for data file: %v", err)
	}
	f.Close()

	index, err := newHashIndexFromFile(filename)
	if err != nil {
		return nil, fmt.Errorf("couldn't build index: %v", err)
	}

	return &HashIndexedFsAppendOnlyStorage{filename: filename, index: index, endOffset: info.Size()}, nil
}

func (s *HashIndexedFsAppendOnlyStorage) Get(key string) (string, bool, error) {
	// Check index for key
	offset, exists := s.index[key]
	if !exists {
		return "", false, nil
	}

	// Open the data file and seek to the relevant line offset
	f, err := os.Open(s.filename)
	if err != nil {
		return "", false, fmt.Errorf("couldn't open data file: %v", err)
	}
	defer f.Close()
	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		return "", false, fmt.Errorf("couldn't seek to offset %d in file: %v", offset, err)
	}

	// Read the line from the file
	r := bufio.NewReader(f)
	line, _, err := r.ReadLine()
	if err != nil {
		return "", false, fmt.Errorf("couldn't read line at offset %d in file: %v", offset, err)
	}

	// Parse key and value from the line and quickly sanity check
	keyValuePair := strings.SplitN(string(line), ", ", 2)
	readKey, value := keyValuePair[0], keyValuePair[1]
	if readKey != key {
		return "", false, fmt.Errorf("key at offset %d is '%s', expected '%s'", offset, readKey, key)
	}

	return value, exists, nil
}

func (s *HashIndexedFsAppendOnlyStorage) Set(key string, value string) error {
	// TODO: keep a file pointer open and don't open and close each time?
	f, err := os.OpenFile(s.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open data file: %v", err)
	}
	defer f.Close()

	recordStartOffset := s.endOffset
	nBytes, err := f.WriteString(key + ", " + value + "\n") // go strings are utf8
	s.endOffset += int64(nBytes)
	if err != nil {
		return fmt.Errorf("couldn't append to data file: %v", err)
	}

	// only save offset once record has already been written to avoid race conditions
	s.index[key] = recordStartOffset

	return nil
}

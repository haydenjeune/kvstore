package store

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type FsAppendOnlyStorage struct {
	filename string
}

func NewFsAppendOnlyStorage(filename string) (*FsAppendOnlyStorage, error) {
	// Ensure data file exists and is openable
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("couldn't open data file: %v", err)
	}
	f.Close()
	return &FsAppendOnlyStorage{filename: filename}, nil
}

func (s *FsAppendOnlyStorage) Get(key string) (string, bool, error) {
	f, err := os.Open(s.filename)
	if errors.Is(err, fs.ErrNotExist) {
		// The case where the storage file does not yet exist is defined as the key not existing
		return "", false, nil
	} else if err != nil {
		return "", false, fmt.Errorf("couldn't open data file: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	value := ""
	exists := false
	var line string
	for scanner.Scan() {
		line = string(scanner.Bytes()) // go strings are utf8
		if strings.HasPrefix(line, key+", ") {
			value = strings.SplitN(line, ", ", 2)[1]
			exists = true
		}
	}
	if scanner.Err() != nil {
		return "", false, fmt.Errorf("couldn't scan data file: %v", scanner.Err())
	}
	return value, exists, nil
}

func (s *FsAppendOnlyStorage) Set(key string, value string) error {
	f, err := os.OpenFile(s.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("couldn't open data file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(key + ", " + value + "\n") // go strings are utf8
	if err != nil {
		return fmt.Errorf("couldn't append to data file: %v", err)
	}

	return nil
}

package store

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FsAppendOnlyStorage struct {
	filename string
}

func NewFsAppendOnlyStorage(filename string) *FsAppendOnlyStorage {
	return &FsAppendOnlyStorage{filename: filename}
}

func (s *FsAppendOnlyStorage) Get(key string) (string, bool, error) {
	f, err := os.Open(s.filename)
	if err != nil {
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
			value = strings.SplitAfterN(line, ", ", 2)[1]
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

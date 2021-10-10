package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	storage := NewFsAppendOnlyStorage("data.kvstore")
	flag.Parse()
	command := flag.Arg(0)
	if command == "get" {
		if flag.NArg() != 2 {
			fmt.Println("Error: expected a single argument for get")
			os.Exit(1)
		}
		key := flag.Arg(1)
		result, err := storage.Get(key)
		if err != nil {
			fmt.Printf("Error: failed to get key %s: %v\n", key, err)
			os.Exit(1)
		}
		fmt.Print(result)

	} else if command == "set" {
		if flag.NArg() != 3 {
			fmt.Println("Error: expected a key and value argument for set")
			os.Exit(1)
		}
		key := flag.Arg(1)
		value := flag.Arg(2)
		err := storage.Set(key, value)
		if err != nil {
			fmt.Printf("Error: failed to get key %s: %v\n", key, err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("Error: first argument must be 'set' or 'get', got '%s'\n", command)
		os.Exit(1)
	}
}

type FsAppendOnlyStorage struct {
	filename string
}

func NewFsAppendOnlyStorage(filename string) *FsAppendOnlyStorage {
	return &FsAppendOnlyStorage{filename: filename}
}

func (s *FsAppendOnlyStorage) Get(key string) (string, error) {
	f, err := os.Open(s.filename)
	if err != nil {
		return "", fmt.Errorf("couldn't open data file: %v", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	value := ""
	var line string
	for scanner.Scan() {
		line = string(scanner.Bytes()) // go strings are utf8
		if strings.HasPrefix(line, key+", ") {
			value = strings.SplitAfterN(line, ", ", 2)[1]
		}
	}
	if scanner.Err() != nil {
		return "", fmt.Errorf("couldn't scan data file: %v", scanner.Err())
	}
	return value, nil
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

package sortedfile

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/haydenjeune/kvstore/pkg/bst"
	"github.com/spf13/afero"
)

type SortedFileKvStorage struct {
	fs       afero.Fs
	memtable bst.BinarySearchTree
	files    []*SortedFile // ordered oldest to newest
}

func NewSortedFileKvStorage(fs afero.Fs) (*SortedFileKvStorage, error) {
	return &SortedFileKvStorage{fs: fs}, nil
}

func (s *SortedFileKvStorage) Get(key string) (string, bool, error) {
	value, exists := s.memtable.Search(key)
	if exists {
		return value, true, nil
	}

	for i := len(s.files) - 1; i >= 0; i-- {
		value, exists, err := s.files[i].Get(key)
		if err != nil {
			return "", false, fmt.Errorf("failed to get key '%s': %v", key, err)
		} else if exists {
			return value, true, nil
		}
	}

	return "", false, nil
}

func (s *SortedFileKvStorage) Set(key string, value string) error {
	s.memtable.Insert(key, value)

	if s.memtable.Size() >= MAX_RECORDS_PER_FILE {
		filename, err := s.nextFileName()
		if err != nil {
			return fmt.Errorf("failed to infer next filename in series: %v", err)
		}
		err = writeBstToSortedFile(&s.memtable, filename, s.fs)
		if err != nil {
			return fmt.Errorf("failed to write memtable to file: %v", err)
		}
		file, err := NewSortedFile(filename, s.fs)
		if err != nil {
			return fmt.Errorf("failed to read new sorted file: %v", err)
		}
		s.files = append(s.files, file)
		s.memtable = bst.BinarySearchTree{}
	}

	return nil
}

func (s *SortedFileKvStorage) nextFileName() (string, error) {
	files, err := afero.ReadDir(s.fs, ".")
	if err != nil {
		return "", fmt.Errorf("failed to list files: %v", err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().UnixMicro() < files[j].ModTime().UnixMicro()
	})

	// all filenames are just integers, starting at 0
	var last int = -1
	for _, f := range files {
		if !f.IsDir() {
			fileNumber, err := strconv.Atoi(f.Name())
			if err == nil {
				last = fileNumber
			}
		}
	}

	return strconv.Itoa(last + 1), nil
}

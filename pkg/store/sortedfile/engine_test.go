package sortedfile

import (
	"path"
	"strconv"
	"testing"

	"github.com/spf13/afero"
)

func Test_nextFileNameForDir_Returns0ForEmptyDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	s, err := NewSortedFileKvStorage(fs)
	if err != nil {
		t.Fatalf("Unexpected error initialising storage: %v", err)
	}

	result, err := s.nextFileName()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "0"
	if path.Base(result) != expected {
		t.Fatalf("Expected filename '%s', got '%s'", expected, result)
	}
}

func Test_nextFileNameForDir_ReturnsCorrectNextFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.Create("0")
	fs.Create("1")
	fs.Create("3") // gaps should be skipped
	s, err := NewSortedFileKvStorage(fs)
	if err != nil {
		t.Fatalf("Unexpected error initialising storage: %v", err)
	}

	result, err := s.nextFileName()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := "4"
	if path.Base(result) != expected {
		t.Fatalf("Expected filename '%s', got '%s'", expected, result)
	}
}

func Test_SortedFileKvStorage_CorrectlyWritesSortedFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	storage, err := NewSortedFileKvStorage(fs)
	if err != nil {
		t.Fatalf("Failed to init storage: %v", err)
	}

	for i := 0; i < 150; i++ {
		err = storage.Set(strconv.Itoa(i), strconv.Itoa(i*2))
		if err != nil {
			t.Fatalf("Failed to set %dth record: %v", i, err)
		}
	}

	result, exists, err := storage.Get("33")
	if err != nil || !exists || result != "66" {
		t.Fatalf("Failed to get record '33': %v", err)
	}

	exists, _ = afero.Exists(fs, "0")
	if !exists {
		t.Fatalf("Expected a file with name '0' to be written")
	}
}

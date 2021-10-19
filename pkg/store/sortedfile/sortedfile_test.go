package sortedfile

import (
	"path"
	"strconv"
	"testing"

	"github.com/spf13/afero"
)

func Test_newSparseIndexFromFile_ReadsCorrectFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	f, err := fs.Create("testfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	f.WriteString(`a, value1
b, value2
c, value3
`)
	f.Close()

	index, err := newSparseIndexFromFile(f.Name(), fs)
	if err != nil {
		t.Fatal("Failed to read storage file")
	}
	expected := []KeyOffset{
		{Key: "a", Offset: 0},
	}

	if index[0] != expected[0] {
		t.Fatal("Unexpected value")
	}
}

func Test_newSparseIndexFromFile_ErrorsForUnorderedFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	f, err := fs.Create("testfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	f.WriteString(`b, value1
a, value2
c, value3
`)
	f.Close()

	_, err = newSparseIndexFromFile(f.Name(), fs)
	if err == nil {
		t.Fatal("Expected error on reading badly ordered file")
	}
}

func Test_newSparseIndexFromFile_ErrorsForDuplicateKeys(t *testing.T) {
	fs := afero.NewMemMapFs()
	f, err := fs.Create("testfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	f.WriteString(`a, value1
a, value2
c, value3
`)
	f.Close()

	_, err = newSparseIndexFromFile(f.Name(), fs)
	if err == nil {
		t.Fatal("Expected error on reading file with duplicate keys")
	}
}

func Test_newSparseIndexFromFile_ReadsEmptyKey(t *testing.T) {
	fs := afero.NewMemMapFs()
	f, err := fs.Create("testfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	f.WriteString(", value1\n")
	f.Close()

	index, err := newSparseIndexFromFile(f.Name(), fs)
	if err != nil {
		t.Fatal("Failed to read file")
	}

	expected := KeyOffset{Key: "", Offset: 0}
	if index[0] != expected {
		t.Fatal("Failed to parse empty key")
	}
}

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

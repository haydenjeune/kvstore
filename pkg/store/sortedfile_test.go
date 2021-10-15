package store

import (
	"os"
	"testing"
)

func Test_newSparseIndexFromFile_ReadsCorrectFile(t *testing.T) {
	f, err := os.CreateTemp("", "tempfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	filename := f.Name()
	defer os.Remove(filename)

	f.WriteString(`a, value1
b, value2
c, value3
`)
	f.Close()

	index, err := newSparseIndexFromFile(filename)
	if err != nil {
		t.Fatal("Failed to read storage file")
	}
	expected := []KeyOffset{
		KeyOffset{Key: "a", Offset: 0},
		KeyOffset{Key: "b", Offset: 10},
		KeyOffset{Key: "c", Offset: 20},
	}

	for i:=0;i<3;i++ {
		if index[i] != expected[i] {
			t.Fatalf("Unexpected value at index %d", i)
		}
	}
}

func Test_newSparseIndexFromFile_ErrorsForUnorderedFile(t *testing.T) {
	f, err := os.CreateTemp("", "tempfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	filename := f.Name()
	defer os.Remove(filename)

	f.WriteString(`b, value1
a, value2
c, value3
`)
	f.Close()

	_, err = newSparseIndexFromFile(filename)
	if err == nil {
		t.Fatal("Expected error on reading badly ordered file")
	}
}

func Test_newSparseIndexFromFile_ErrorsForDuplicateKeys(t *testing.T) {
	f, err := os.CreateTemp("", "tempfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	filename := f.Name()
	defer os.Remove(filename)

	f.WriteString(`a, value1
a, value2
c, value3
`)
	f.Close()

	_, err = newSparseIndexFromFile(filename)
	if err == nil {
		t.Fatal("Expected error on reading file with duplicate keys")
	}
}

func Test_newSparseIndexFromFile_ReadsEmptyKey(t *testing.T) {
	f, err := os.CreateTemp("", "tempfile")
	if err != nil {
		t.Fatal("Failed to create temporary file")
	}
	filename := f.Name()
	defer os.Remove(filename)

	f.WriteString(", value1\n")
	f.Close()

	index, err := newSparseIndexFromFile(filename)
	if err != nil {
		t.Fatal("Failed to read file")
	}

	expected := KeyOffset{Key: "", Offset:0}
	if index[0] != expected {
		t.Fatal("Failed to parse empty key")
	}
}
package sortedfile

import (
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

func Test_getInterval_EmptySlice(t *testing.T) {
	data := []KeyOffset{}

	var l, r *KeyOffset

	l, r = getInterval(data, "b")
	if l != nil || r != nil {
		t.Fail()
	}
}

func Test_getInterval_LengthOne(t *testing.T) {
	data := []KeyOffset{
		{Key: "f"},
	}

	var l, r *KeyOffset

	l, r = getInterval(data, "b")
	if l != nil || r != &data[0] {
		t.Fail()
	}

	l, r = getInterval(data, "g")
	if l != &data[0] || r != nil {
		t.Fail()
	}

	l, r = getInterval(data, "f")
	if l != &data[0] || r != nil {
		t.Fail()
	}
}

func Test_getInterval(t *testing.T) {
	data := []KeyOffset{
		{Key: "b"},
		{Key: "e"},
		{Key: "f"},
		{Key: "k"},
		{Key: "s"},
	}

	var l, r *KeyOffset

	l, r = getInterval(data, "a")
	if l != nil || r != &data[0] {
		t.Fail()
	}

	l, r = getInterval(data, "b")
	if l != &data[0] || r != &data[1] {
		t.Fail()
	}

	l, r = getInterval(data, "c")
	if l != &data[0] || r != &data[1] {
		t.Fail()
	}

	l, r = getInterval(data, "e")
	if l != &data[1] || r != &data[2] {
		t.Fail()
	}

	l, r = getInterval(data, "f")
	if l != &data[2] || r != &data[3] {
		t.Fail()
	}

	l, r = getInterval(data, "m")
	if l != &data[3] || r != &data[4] {
		t.Fail()
	}

	l, r = getInterval(data, "z")
	if l != &data[4] || r != nil {
		t.Fail()
	}
}

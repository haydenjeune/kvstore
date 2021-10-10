package bst

import (
	"testing"
)

func Test_Search_AfterInsert_ReturnsExpected(t *testing.T) {
	key := "test"
	value := "blahblah"
	sut := BinaryNode{}
	sut.Insert(key, value)

	result, exists := sut.Search(key)

	if exists != true {
		t.Error("Search returned key does not exist when it should")
	}
	if result != value {
		t.Errorf("Search returned value '%s' but should have returned '%s'", result, value)
	}
}

func Test_Search_AfterMultipleInsert_ReturnsLatestValue(t *testing.T) {
	key := "test"
	value := "blahblah"
	sut := BinaryNode{}
	sut.Insert(key, "notwhatwewant")
	sut.Insert(key, value)

	result, exists := sut.Search(key)

	if exists != true {
		t.Error("Search returned key does not exist when it should")
	}
	if result != value {
		t.Errorf("Search returned value '%s' but should have returned '%s'", result, value)
	}
}

func Test_Search_AfterInsertWithDifferentKey_ReturnsNoValue(t *testing.T) {
	sut := BinaryNode{}
	sut.Insert("adifferentkey", "notwhatwewant")

	_, exists := sut.Search("key")

	if exists != false {
		t.Error("Search returned key exists when it should not")
	}
}

func Test_InOrderTraversal_ReturnsCorrectNumberOfElements(t *testing.T) {
	tree := BinarySearchTree{}
	tree.Insert("5", "five")
	tree.Insert("4", "four")
	tree.Insert("3", "three")

	iter := NewInOrderTraversalIterator(&tree)
	counter := 0

	for iter.Next() != nil {
		counter += 1
	}

	if counter != 3 {
		t.Fatalf("Put 3 nodes in, but the iterator returned %d", counter)
	}
}

func Test_InOrderTraversal_ReturnsSortedList(t *testing.T) {
	tree := BinarySearchTree{}
	tree.Insert("5", "five")
	tree.Insert("4", "four")
	tree.Insert("3", "three")
	tree.Insert("6", "six")
	tree.Insert("8", "eight")
	tree.Insert("7", "seven")

	iter := NewInOrderTraversalIterator(&tree)

	for i, key := range []string{"3", "4", "5", "6", "7", "8"} {
		node := iter.Next()
		if node.key != key {
			t.Fatalf("Bad key ordering and %dth element. Got '%s', expected '%s'", i, node.key, key)
		}
	}
}

func Test_BinarySearchTree(t *testing.T) {
	tree := BinarySearchTree{}

	if tree.Size() != 0 {
		t.Fail()
	}

	tree.Insert("1", "one")
	tree.Insert("3", "three")
	tree.Insert("key", "WOW")

	if result, _ := tree.Search("1"); result != "one" {
		t.Fail()
	}
	if result, _ := tree.Search("3"); result != "three" {
		t.Fail()
	}
	if result, _ := tree.Search("key"); result != "WOW" {
		t.Fail()
	}
	if _, exists := tree.Search("DoesNotExist"); exists != false {
		t.Fail()
	}
	if tree.Size() != 3 {
		t.Fail()
	}
}

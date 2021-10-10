package main

import "testing"

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
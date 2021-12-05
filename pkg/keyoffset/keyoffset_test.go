package keyoffset

import "testing"

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



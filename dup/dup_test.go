package dup

import "testing"

func TestDup(t *testing.T) {
	d := NewDup()

	id := d.Random()
	exist := d.Check(id)

	if exist {
		t.Fail()
	}

	d.Track(id)

	exist = d.Check(id)

	if !exist {
		t.Fail()
	}
}

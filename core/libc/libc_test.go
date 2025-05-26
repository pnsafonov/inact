package libc

import "testing"

func TestGetutent0(t *testing.T) {
	ents := Getutent0()
	l0 := len(ents)
	if l0 != 0 {
		t.Fatalf("Getutent0() returned %d entries", l0)
	}
}

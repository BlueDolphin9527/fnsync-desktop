package clip

import (
	"testing"
)

func TestClip1(t *testing.T) {
	pb = nil
	ok := Set("test1")
	if pb == nil {
		t.Errorf(`1: Set() failed, pb == nil\n`)
	}
	if !ok {
		t.Errorf(`2: Set() failed\n`)
	}
	x := Get()
	if x != "test1" {
		t.Errorf(`3: expected "test1", got "%s"\n`, x)
	}
	pb = nil
	Get() // should not panic
	Clear()
	x = Get()
	if x != "" {
		t.Errorf(`4: pasteboard not cleared, got "%s"\n`, x)
	}
	ok = Set("test2")
	if !ok {
		t.Errorf(`5: Set() failed\n`)
	}
	x = Get()
	if x != "test2" {
		t.Errorf(`6: expected "test2", got "%s"\n`, x)
	}
	Clear()
	x = Get()
	if x != "" {
		t.Errorf(`7: pasteboard was not cleared, got "%s"\n`, x)
	}

	pb = nil
	ok = Set("test3")
	if pb == nil {
		t.Errorf(`8: Set() failed, pb == nil\n`)
	}
	if !ok {
		t.Errorf("9: Set() failed\n")
	}
	x = Get()
	if x != "test3" {
		t.Errorf(`10: expected "test3", got "%s"\n`, x)
	}
}

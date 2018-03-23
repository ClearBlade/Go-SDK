package uuid

import (
	"reflect"
	"testing"
)

func TestReality(t *testing.T) {
	if true == false {
		t.Error("Goodbye Cruel World")
	}
}

func TestToAndFromString(t *testing.T) {
	u, _ := New()
	newU, _ := ParseUuid(u.String())
	if !reflect.DeepEqual(u, newU) {
		t.Error("Parse Uuid doesn't work")
	}
}

package services

import (
	"testing"
	"time"
)

func TestGetKeySetKey(t *testing.T) {
	key := "abc"
	val := "def"
	SetKey(key, val, 1000*time.Millisecond)
	newVal, err := GetKey(key)
	if err != nil {
		t.Error(err)
	}
	if newVal != val {
		t.Error(err)
	}
}

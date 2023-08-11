package utils

import (
	"bytes"
	"testing"
)

func TestHash256(t *testing.T) {
	print(bytes.Equal(Hash256("1000"), Hash256("1000")))
}

package dktest

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	sizes := []uint{0, 1, 2, 10, 100}

	for _, s := range sizes {
		t.Run(fmt.Sprintf("size %d", s), func(t *testing.T) {
			str := randString(s)
			if uint(len(str)) != s {
				t.Error("Got wrong randString size:", len(str), "!=", s)
			}
		})
	}
}

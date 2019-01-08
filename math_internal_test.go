package dktest

import (
	"fmt"
	"testing"
)

func TestMin(t *testing.T) {
	testCases := []struct {
		a           int
		b           int
		expectedMin int
	}{
		// both positive
		{a: 1, b: 2, expectedMin: 1},
		{a: 2, b: 1, expectedMin: 1},
		// both negative
		{a: -1, b: -2, expectedMin: -2},
		{a: -2, b: -1, expectedMin: -2},
		// negative and positive
		{a: -1, b: 1, expectedMin: -1},
		{a: 1, b: -1, expectedMin: -1},
		// equal
		{a: 0, b: 0, expectedMin: 0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("min(%d, %d)", tc.a, tc.b), func(t *testing.T) {
			if m := min(tc.a, tc.b); m != tc.expectedMin {
				t.Error("min doesn't match expected:", m, "!=", tc.expectedMin)
			}
		})
	}
}

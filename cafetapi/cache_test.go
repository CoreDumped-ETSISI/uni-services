package main

import (
	"testing"
	"time"
)

func TestIsCacheValid(t *testing.T) {
	testcases := [][]int{
		{1, 5, 1, 1}, // from, to, day, expected
		{1, 5, 2, 1},
		{1, 5, 5, 1},
		{26, 2, 26, 1},
		{26, 2, 30, 1},
		{26, 2, 1, 1},
		{26, 2, 2, 1},

		{1, 5, 6, 0},
		{26, 2, 15, 0},
		{0, 0, 12, 1},
	}

	for _, test := range testcases {
		expected := test[3] == 1
		if expected != isCacheValid(test[0], test[1], time.Date(2018, 1, test[2], 0, 0, 0, 0, time.Local)) {
			t.Logf("FAILED: %v in [%v, %v], expected: %v, got: %v", test[2], test[0], test[1], expected, !expected)
			t.Fail()
		}
	}
}

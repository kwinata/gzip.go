package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRLE(t *testing.T) {
	values := []int{
		3, 0, 0, 0, 4, // 0-4
		4, 3, 2, 3, 3, // 5-9
		4, 5, 0, 0, 0, // 10-14
		0, 6, 7, 7, // 15-18
	}
	expectedRanges := []rleRange{
		{0, 3},
		{3, 0},
		{5, 4},
		{6, 3},
		{7, 2},
		{9, 3},
		{10, 4},
		{11, 5},
		{15, 0},
		{16, 6},
		{18, 7},
	}
	assert.Equal(t, expectedRanges, runLengthEncoding(values))
}

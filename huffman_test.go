package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildHuffmanTree(t *testing.T) {
	testCases := []struct {
		hRanges              []huffmanRange
		expectedHuffmanCodes []string
	}{
		{
			[]huffmanRange{
				{1, 4},
				{4, 6},
				{6, 4},
				{14, 5},
				{18, 6},
				{21, 4},
				{26, 6},
			},
			[]string{
				"0000",
				"0001",
				"101100",
				"101101",
				"101110",
				"0010",
				"0011",
				"01110",
				"01111",
				"10000",
				"10001",
				"10010",
				"10011",
				"10100",
				"10101",
				"101111",
				"110000",
				"110001",
				"110010",
				"0100",
				"0101",
				"0110",
				"110011",
				"110100",
				"110101",
				"110110",
				"110111",
			},
		},
		{
			[]huffmanRange{
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
			},
			[]string{
				// total of 19 codes

				// 0
				"010",
				"",
				"",
				"",
				"1100",

				// 5
				"1101",
				"011",
				"00",
				"100",
				"101",

				// 10
				"1110",
				"11110",
				"",
				"",
				"",

				// 15
				"",
				"111110",
				"1111110",
				"1111111",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test %d", i+1), func(t *testing.T) {
			root := buildHuffmanTree(tc.hRanges)
			codeTable := make([]string, tc.hRanges[len(tc.hRanges)-1].end+1) // symbol 0 is valid
			traverseHuffmanTree(root, "", codeTable)
			for i, v := range codeTable {
				assert.Equal(t, tc.expectedHuffmanCodes[i], v)
			}
		})
	}

}

func TestBuildHRangesFromBitLengthsArray(t *testing.T) {
	bitLengths := []int{
		3, 0, 0, 0, 4, // 0-4
		4, 3, 2, 3, 3, // 5-9
		4, 5, 0, 0, 0, // 10-14
		0, 6, 7, 7, // 15-18
	}
	expectedHRanges := []huffmanRange{
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
	assert.Equal(t, expectedHRanges, runLengthEncoding(bitLengths))
}

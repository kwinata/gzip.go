package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadCodes(t *testing.T) {
	hclen := 8 // there are twelve codes
	/*
		110	111	111	011	011	010	011	011	100	100	101	100
		16	17	18	0 	8 	7 	9 	6 	10	5 	11	4
		6	7	7	3	3	2	3	3	4	4	5	4
	*/
	stream := &bitstream{
		source: bytes.NewReader([]byte{
			0b11_111_110,
			0b0_011_011_1,
			0b011_011_01,
			0b01_100_100,
			0b100_1,
		}),
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
	hRanges := readCodes(stream, hclen)
	assert.Equal(t, expectedHRanges, hRanges)

	expectedCodes := []string{
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
	}
	codeTable := make([]string, 19)
	traverseHuffmanTree(buildHuffmanTree(hRanges), "", codeTable)
	assert.Equal(t, expectedCodes, codeTable)
}

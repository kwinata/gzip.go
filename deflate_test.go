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

	expectedBitLengths := []int{
		3, 0, 0, 0, 4, // 0-4
		4, 3, 2, 3, 3, // 5-9
		4, 5, 0, 0, 0, // 10-14
		0, 6, 7, 7, // 15-18
	}
	codesBitLengths := readCodesBitLengths(stream, hclen)
	assert.Equal(t, expectedBitLengths, codesBitLengths)
}

func TestReadAlphabetsBitLengths(t *testing.T) {
	codesHRanges := []huffmanRange{
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
	codesHuffmanTreeRoot := buildHuffmanTree(codesHRanges)
	debugCodeTable := make([]string, 19)
	traverseHuffmanTree(codesHuffmanTreeRoot, "", debugCodeTable)

	source := bytes.NewReader(helperBitStringToBytes("1111110111001111111010100011011011001110"))
	/*
		1111110	 111
		17	     repeat '0' 10

		00
		7 (literal)

		1111111	0101000
		18	    repeat '0' 21 times

		1101
		5 (literal)

		101
		9 (literal)

		100
		8 (literal)

		1110
		10 (literal)
	*/
	alphabetBitLengths := []int{
		// there are 10 + 1 (7) + 21 + 4 (5, 9, 8, 10) = 36 alphabets
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		7, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 5, 9, 8, 10,
	}
	assert.Equal(t, alphabetBitLengths,
		readAlphabetsBitLengths(&bitstream{source: source}, len(alphabetBitLengths), codesHuffmanTreeRoot))
}

func TestInflateHuffmanCodesNoBackPointer(t *testing.T) {
	// These are inefficient huffman trees. This is used to make it easier to create the test cases
	literalsRoot := buildHuffmanTree([]huffmanRange{
		{285, 16},
	})
	distancesRoot := buildHuffmanTree([]huffmanRange{
		{30, 8},
	})

	stream := &bitstream{
		source: bytes.NewReader([]byte{
			0x00, 0x00, // 0x00
			0x00, 0x80, // 0x01
			0x00, 0x40, // 0x02
			0x00, 0x20, // 0x04
			0x80, 0x00, // stop code
		}),
	}
	outBytes := inflateHuffmanCodes(stream, literalsRoot, distancesRoot)
	assert.Equal(t, []byte{0x00, 0x01, 0x02, 0x04}, outBytes)
}

func TestInflateHuffmanCodesWithLiteralBackPointer(t *testing.T) {
	// These are inefficient huffman trees. This is used to make it easier to create the test cases
	literalsRoot := buildHuffmanTree([]huffmanRange{
		{285, 16},
	})
	distancesRoot := buildHuffmanTree([]huffmanRange{
		{30, 8},
	})

	stream := &bitstream{
		source: bytes.NewReader([]byte{
			0x00, 0x00, // 0x00
			0x00, 0x80, // 0x01
			0x00, 0x40, // 0x02
			0x00, 0x20, // 0x04
			0x80, 0x40, // Code 258, back-pointer length of 4
			0x40, // distance code 2,
			0x00, 0xC0, // 0x03
			0x80, 0x00, // stop code
		}),
	}
	outBytes := inflateHuffmanCodes(stream, literalsRoot, distancesRoot)
	assert.Equal(t, []byte{0x00, 0x01, 0x02, 0x04, 0x01, 0x02, 0x04, 0x01, 0x03}, outBytes)
}

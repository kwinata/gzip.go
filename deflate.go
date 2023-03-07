package main

import (
	"fmt"
	"time"
)

func readFixedHuffmanTree(stream *bitstream) (root *huffmanNode) {
	return buildHuffmanTree([]rleRange{
		{143, 8},
		{255, 9},
		{279, 7},
		{287, 8},
	})
}

func readDynamicHuffmanTree(stream *bitstream) (literalsRoot *huffmanNode, distancesRoot *huffmanNode) {
	/*
		format is:
		- header (hlit|hdist|hclen)
		- (hclen+4) * 3 code lengths
		- followed by 257 + hdist + hlit literals

		hdist and hlit should help to define the distance and length codes
	*/

	var hlit = readBitsInv(stream, 5)
	var hdist = readBitsInv(stream, 5)

	// there are (hclen + 4) number of codes
	var hclen = readBitsInv(stream, 4)

	if explanationMode {
		fmt.Printf("hlit: %d (number of (extra) length literals)\n", hlit)
		fmt.Printf("hdist: %d (number of distance codes)\n", hdist)
		fmt.Printf("hclen: %d (number of huffman code length for the first tree)\n", hclen)
	}

	// read codes
	codeBitLengths := readCodesBitLengths(stream, hclen)
	codeHuffmanRoot := buildHuffmanTree(runLengthEncoding(codeBitLengths))

	// read alphabet
	alphabetsBitLengths := readAlphabetsBitLengths(stream, 258+hlit+hdist, codeHuffmanRoot)

	alphabetsBitLengths = append(alphabetsBitLengths)

	// split alphabets into literals and distances
	literals := alphabetsBitLengths[:hlit+257]
	distances := alphabetsBitLengths[hlit+257:]
	literalsRLE := runLengthEncoding(append([]int{0}, literals...)) // Seems to be using 1-indexing
	distancesRLE := runLengthEncoding(distances)
	literalsRoot = buildHuffmanTree(literalsRLE)
	distancesRoot = buildHuffmanTree(distancesRLE)

	return
}

func readCodesBitLengths(stream *bitstream, hclen int) []int {
	// The specification refers to the repetition codes as the values 16, 17, and 18,
	//	but these numbers don't have any real physical meaning
	//	16 means "repeat the previous character n times",
	//	and 17 & 18 mean "insert n 0's".
	//	The number n follows the repeat codes and is encoded
	//	(without compression) in 2, 3 or 7 bits, respectively

	// Because codes of lengths 15, 1, 14 and 2 are likely to be very rare in real-world data,
	// 	the codes themselves are given in order of expected frequency
	codeLengthOffsets := []int{
		16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15,
	}

	codeBitLengths := make([]int, 19) // max hclen (0b1111) + 4
	for i := 0; i < (hclen + 4); i++ {
		codeBitLengths[codeLengthOffsets[i]] = readBitsInv(stream, 3)
	}

	return codeBitLengths
}

func readAlphabetsBitLengths(stream *bitstream, alphabetCount int, codeLengthsRoot *huffmanNode) []int {
	alphabetBitLengths := make([]int, alphabetCount)

	i := 0
	for i < alphabetCount {
		code, _ := getCode(stream, codeLengthsRoot)
		// 0-15: literal (4 bits)
		// 16: repeat the previous character n times
		// 17: insert n 0's (3 bit specified), max value is 10
		// 18: insert n 0's (7 bit specifier), add 11 (because it's the max of code 17)
		if code == 16 {
			repeatLength := readBitsInv(stream, 2) + 3
			for j := 0; j < repeatLength; j++ {
				alphabetBitLengths[i] = alphabetBitLengths[i-1]
				i++
			}
		} else if code == 17 || code == 18 {
			var repeatLength int
			if code == 17 {
				repeatLength = readBitsInv(stream, 3) + 3
			} else {
				repeatLength = readBitsInv(stream, 7) + 11
			}
			for j := 0; j < repeatLength; j++ {
				alphabetBitLengths[i] = 0
				i++
			}
		} else {
			alphabetBitLengths[i] = code
			i++
		}
	}
	return alphabetBitLengths
}

/*

reading LZ77:

Format is always Length|Distance

- codes 257-264: length is $code - 254 (no extra length bits)
- codes 265-285: have extra bits length


Distance codes:
Can range from 1-32768
*/

var shouldPrintInline = true

func inflateHuffmanCodes(stream *bitstream, literalsRoot *huffmanNode, distancesRoot *huffmanNode) []byte {
	/*
		Now, if there are only 285-257=28 length codes, that doesn't give the LZ77 compressor much room to
		reuse previous input. Instead, the deflate format uses the 28 pointer codes as an indication to the
		decompressor as to how many extra bits follow which indicate the actual length of the match.
	*/

	/*
		What's with this extraLengthAddend?
		It is used as: length = readBitsInv(stream, (node.code - 261)/4) + extraLengthAddend[node.code - 265]
		for node.code in [265, 285)
		Code: 265, base value: 11, max_value: 12
		Code: 266, base value: 13, max_value: 14
		Code: 267, base value: 15, max_value: 16
		Code: 268, base value: 17, max_value: 18
		Code: 269, base value: 19, max_value: 22
		Code: 270, base value: 23, max_value: 26
		Code: 271, base value: 27, max_value: 30
		Code: 272, base value: 31, max_value: 34
		Code: 273, base value: 35, max_value: 42
		Code: 274, base value: 43, max_value: 50
		Code: 275, base value: 51, max_value: 58
		Code: 276, base value: 59, max_value: 66
		Code: 277, base value: 67, max_value: 82
		Code: 278, base value: 83, max_value: 98
		Code: 279, base value: 99, max_value: 114
		Code: 280, base value: 115, max_value: 130
		Code: 281, base value: 131, max_value: 162
		Code: 282, base value: 163, max_value: 194
		Code: 283, base value: 195, max_value: 226
		Code: 284, base value: 227, max_value: 258
	*/
	extraLengthAddend := []int{
		11, 13, 15, 17, 19, 23, 27,
		31, 35, 43, 51, 59, 67, 83,
		99, 115, 131, 163, 195, 227,
	}

	/*
		We only support until distance code 29 instead of until 31 because it's sufficient to describe until 32KiB distance
		Dist Code: 4, base value: 4, max_value: 5
		Dist Code: 5, base value: 6, max_value: 7
		Dist Code: 6, base value: 8, max_value: 11
		Dist Code: 7, base value: 12, max_value: 15
		Dist Code: 8, base value: 16, max_value: 23
		Dist Code: 9, base value: 24, max_value: 31
		Dist Code: 10, base value: 32, max_value: 47
		Dist Code: 11, base value: 48, max_value: 63
		Dist Code: 12, base value: 64, max_value: 95
		Dist Code: 13, base value: 96, max_value: 127
		Dist Code: 14, base value: 128, max_value: 191
		Dist Code: 15, base value: 192, max_value: 255
		Dist Code: 16, base value: 256, max_value: 383
		Dist Code: 17, base value: 384, max_value: 511
		Dist Code: 18, base value: 512, max_value: 767
		Dist Code: 19, base value: 768, max_value: 1023
		Dist Code: 20, base value: 1024, max_value: 1535
		Dist Code: 21, base value: 1536, max_value: 2047
		Dist Code: 22, base value: 2048, max_value: 3071
		Dist Code: 23, base value: 3072, max_value: 4095
		Dist Code: 24, base value: 4096, max_value: 6143
		Dist Code: 25, base value: 6144, max_value: 8191
		Dist Code: 26, base value: 8192, max_value: 12287
		Dist Code: 27, base value: 12288, max_value: 16383
		Dist Code: 28, base value: 16384, max_value: 24575
		Dist Code: 29, base value: 24576, max_value: 32767
	*/
	extraDistAddend := []int{
		4, 6, 8, 12, 16, 24, 32, 48,
		64, 96, 128, 192, 256, 384,
		512, 768, 1024, 1536, 2048,
		3072, 4096, 6144, 8192,
		12288, 16384, 24576,
	}
	node := literalsRoot
	buf := make([]byte, 0)
	var debugNode []byte
	for {
		if nextBit(stream) != 0 {
			node = node.one
			debugNode = append(debugNode, '1')
		} else {
			node = node.zero
			debugNode = append(debugNode, '0')
		}
		if node == nil {
			panic(fmt.Errorf("%s is not a valid huffman code / path", debugNode))
		}
		if node.code != -1 {
			if shouldPrintInline && slowPrintMode {
				time.Sleep(50 * time.Millisecond)
			}
			node = &huffmanNode{code: node.code - 1}
			debugNode = nil
			if node.code >= 0 && node.code < 256 {
				// literal code
				buf = append(buf, byte(node.code))
				if shouldPrintInline {
					fmt.Printf("%s", string(rune(node.code)))
				}
			} else if node.code == 256 {
				// stop code
				break
			} else if node.code > 256 && node.code <= 285 {
				// This is a back-pointer

				// get length
				var length int
				if node.code < 265 {
					length = node.code - 254
				} else if node.code == 285 {
					length = 258 // this seems to be for a short cut for the 284? not sure why don't we use 259 instead?
				} else {
					length = extraLengthAddend[node.code-265] + readBitsInv(stream, (node.code-261)/4)
				}

				var dist int
				if distancesRoot == nil {
					// hardcoded distances
					dist = readBitsInv(stream, 5)
				} else {
					// get bits (5 bits)
					distanceNode := distancesRoot
					for distanceNode.code == -1 {
						if nextBit(stream) != 0 {
							distanceNode = distanceNode.one
						} else {
							distanceNode = distanceNode.zero
						}
					}
					dist = distanceNode.code
					if dist > 3 {
						extraDist := readBitsInv(stream, (dist-2)/2)
						dist = extraDist + extraDistAddend[dist-4]
					}
				}
				backPointer := len(buf) - dist - 1
				if shouldPrintInline && backPointerMode {
					fmt.Printf("<%d,%d>(", backPointer, length)
				}
				bufString := string(buf)
				if bufString == "" {

				}
				for length > 0 {
					buf = append(buf, buf[backPointer])
					if shouldPrintInline {
						fmt.Printf("%s", string(rune(buf[backPointer])))
					}
					length--
					backPointer++
				}
				if shouldPrintInline && backPointerMode {
					fmt.Printf(")")
				}
			} else {
				panic("invalid code!")
			}
			node = literalsRoot
		}
	}
	return buf
}

package main

func readHuffmanTree(stream *bitstream) {
	var hlit = readBitsInv(stream, 5)
	var hdist = readBitsInv(stream, 5)

	// there are (hclen + 4) number of codes
	var hclen = readBitsInv(stream, 4)

	println(hlit, hdist, hclen)
	/*

		format is:
		- header (hlit|hdist|hclen)
		- (hclen+4) * 3 code lengths
		- followed by 257 + hdist + hlit literals

		hdist and hlit should help to define the distance and length codes
	*/
	//codeRanges := readCodesHRanges(stream, hclen)
	//codeLengthsRoot := buildHuffmanTree(codeRanges)
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
		code, huffmanCodeString := getCode(stream, codeLengthsRoot)
		println(huffmanCodeString, code)
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

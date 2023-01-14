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
	//codeRanges := readCodes(stream, hclen)
	//codeLengthsRoot := buildHuffmanTree(codeRanges)
}

func readCodes(stream *bitstream, hclen int) []huffmanRange {
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

	codeLengths := make([]int, 19) // max hclen (0b1111) + 4
	for i := 0; i < (hclen + 4); i++ {
		codeLengths[codeLengthOffsets[i]] = readBitsInv(stream, 3)
	}
	hRanges := make([]huffmanRange, 19)
	j := 0
	for i := 0; i < 19; i++ {
		if i > 0 && codeLengths[i] != codeLengths[i-1] {
			j++
		}
		hRanges[j].end = i
		hRanges[j].bitLength = codeLengths[i]
	}
	return hRanges[:j+1]
}

package main

type huffmanNode struct {
	code int // -1 for non-leaf nodes
	zero *huffmanNode
	one *huffmanNode
}

type huffmanRange struct {
	end int  // set bitLength until index end
	bitLength int
}

type treeNode struct {
	len int
	code int
}

func buildHuffmanTree(hRanges []huffmanRange) *huffmanNode {
	// 1. find max bit length
	maxBitLength := 0
	for _, hRange := range hRanges {
		if hRange.bitLength > maxBitLength {
			maxBitLength = hRange.bitLength
		}
	}

	// (2. allocate space (skipped, will allocate near the logic itself))

	// 3. Determine the number of codes for each bit-length
	blCount := map[int]int{}  // number of codes for each bit length
	previousEnd := -1
	for _, hRange := range hRanges {
		if hRange.end - previousEnd <= 0 {
			panic("The end of each huffmanRange must be strictly increasing")
		}
		count := hRange.end - previousEnd
		blCount[hRange.bitLength] += count
		previousEnd = hRange.end
	}

	// 4. Populate nextCode (the starting `code` for each bit length group)
	nextCode := map[int]int{} // the starting 'code' for each bit length group
	code := 0
	for bitLength := 1;  bitLength <= maxBitLength; bitLength++ {
		// the previous starting point was code, we increment it by the bitLengthCount of
		//   the previous group, then we right shift it by 1.
		code = (code + blCount[bitLength - 1]) << 1;

		nextCode[bitLength] = code

	}

	// 5. assign codes to each symbol in range
	numberOfCodes := hRanges[len(hRanges) - 1].end
	tree := make([]treeNode, numberOfCodes+1) // symbol start from zero
	hRangeIdx := 0
	for ti := 0; ti <= numberOfCodes; ti++ { // ti for tree index
		hRange := hRanges[hRangeIdx]
		if ti > hRange.end {
			// move to the next range
			hRangeIdx++
			hRange = hRanges[hRangeIdx]
		}

		tree[ti].len = hRange.bitLength
		tree[ti].code = nextCode[tree[ti].len]
		nextCode[tree[ti].len]++
	}

	// 6. build huffman tree
	root := &huffmanNode{code: -1}
	for ti := 0; ti <= numberOfCodes; ti++ {
		var node *huffmanNode
		node = root
		if tree[ti].len == 0 {
			continue
		}
		// traverse the tree, build node if not exist; bi is bitIndex
		for bi := tree[ti].len; bi > 0; bi-- {
			if (tree[ti].code & (1 << (bi -1))) > 0 { // if the bi-th bit is set
				if node.one == nil {
					node.one = &huffmanNode{code: -1}
				}
				node = node.one
			} else {
				if node.zero == nil {
					node.zero = &huffmanNode{code: -1}
				}
				node = node.zero
			}
		}
		if node.code != -1 {
			panic("this node shouldn't be set before")
		}
		node.code = ti
	}
	return root
}
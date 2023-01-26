package main


type rleRange struct {
	end       int // set bitLength until index end
	bitLength int
}

func runLengthEncoding(values []int) (ranges []rleRange) {
	for i := 0; i < len(values); i++ {
		if i == 0 || values[i] != values[i-1] {
			ranges = append(ranges, rleRange{i, values[i]})
		} else {
			// if bitLength is the same as previous, simply increase the end pointer
			if ranges == nil {
				panic("shouldn't be nil")
			}
			ranges[len(ranges)-1].end = i
		}
	}
	return ranges
}
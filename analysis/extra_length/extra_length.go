package main

import "fmt"

var extraLengthAddend = []int{
	11, 13, 15, 17, 19, 23, 27,
	31, 35, 43, 51, 59, 67, 83,
	99, 115, 131, 163, 195, 227,
}

func main() {
	for code := 265; code < 285; code++ {
		extraBitsLength := (code - 261)/4
		extraLength := extraLengthAddend[code-265]
		fmt.Printf("Code: %d, base value: %d, max_value: %d\n", code, extraLength, extraLength + (1 << extraBitsLength) - 1)
	}
}

/*
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
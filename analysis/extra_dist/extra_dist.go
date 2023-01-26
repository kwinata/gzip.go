package main

import "fmt"

var extraDistAddends = []int{
	4, 6, 8, 12, 16, 24, 32, 48,
	64, 96, 128, 192, 256, 384,
	512, 768, 1024, 1536, 2048,
	3072, 4096, 6144, 8192,
	12288, 16384, 24576,
}

func main() {
	// extra dist addends only have 26 items instead of 28 items (for it to support 32 distance codes 🤔)
	// why don't we support higher distance?
	// Oh this is because, this way, we can support until 32768 distance, which is exactly 2**15, 32KiB of data
	for dist := 4; dist < 30; dist++ {
		extraDistLength := (dist - 2) / 2
		extraDist := extraDistAddends[dist-4]
		fmt.Printf("Dist Code: %d, base value: %d, max_value: %d\n", dist, extraDist, extraDist+(1<<extraDistLength)-1)
	}
}

/*
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

package main

import (
	"fmt"
)

func longestEvenRun(nums []int) (int, int, int) {
	count, maxCount, end := 0, 0, 0
	for i, n := range nums {
		if n%2 == 0 {
			count++
		} else {
			count = 0
		}
		if count > maxCount {
			maxCount = count
			end = i
		}
	}
	if maxCount <= 2 {
		return -1, 0, -1 // no qualifying subarray found
	}
	return end - maxCount + 1, maxCount, end
}

func main() {
	nums := []int{1, 2, 4, 7, 8, 6, 9, 8, 2, 6, 10, 3, 1, 4, 3, 3}

	start, length, end := longestEvenRun(nums)

	if length == 0 {
		fmt.Println("no even run longer than 2")
		return
	}

	fmt.Printf("Start: %d, Length: %d, End: %d\n", start, length, end)
	fmt.Println(nums[end-length+1 : end+1])
}

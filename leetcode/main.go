package main

import (
	"fmt"
	"sync"
)

// chunkResult summarizes one subslice nums[lo : lo+n].
//
// Besides the best run fully inside the chunk, it records how many even
// numbers touch its left and right edges, so runs that cross a boundary
// between two chunks can be rebuilt when the chunks are merged.
type chunkResult struct {
	lo, n     int // position of the chunk in the full slice, and its length
	prefix    int // number of evens at the start of the chunk
	suffix    int // number of evens at the end of the chunk
	bestStart int // start (index in the full slice) of the longest run in the chunk
	bestLen   int // length of that run
}

// processChunk scans one subslice. lo is the index of chunk[0] in the full
// slice, so bestStart is reported as a global index.
func processChunk(chunk []int, lo int) chunkResult {
	r := chunkResult{lo: lo, n: len(chunk)}
	count := 0
	seenOdd := false
	for i, v := range chunk {
		if v%2 == 0 {
			count++
		} else {
			if !seenOdd {
				r.prefix = count // evens before the first odd number
				seenOdd = true
			}
			count = 0
		}
		if count > r.bestLen {
			r.bestLen = count
			r.bestStart = lo + i - count + 1
		}
	}
	if !seenOdd {
		r.prefix = count // the whole chunk is even
	}
	r.suffix = count // evens after the last odd number
	return r
}

// merge combines two adjacent chunks, where a comes directly before b.
func merge(a, b chunkResult) chunkResult {
	r := chunkResult{lo: a.lo, n: a.n + b.n, bestStart: a.bestStart, bestLen: a.bestLen}

	// The run crossing the boundary: evens at the end of a + evens at the start of b.
	if cross := a.suffix + b.prefix; cross > r.bestLen {
		r.bestLen = cross
		r.bestStart = a.lo + a.n - a.suffix
	}
	// Strict > keeps the earliest run on ties, same as the sequential version.
	if b.bestLen > r.bestLen {
		r.bestLen = b.bestLen
		r.bestStart = b.bestStart
	}

	// If a is all even, its run continues into b's prefix (and vice versa for
	// b's suffix), so the merged chunk's edges can grow across the boundary.
	r.prefix = a.prefix
	if a.prefix == a.n {
		r.prefix = a.n + b.prefix
	}
	r.suffix = b.suffix
	if b.suffix == b.n {
		r.suffix = b.n + a.suffix
	}
	return r
}

// longestEvenRunParallel splits nums into parts subslices, scans each in its
// own goroutine, then merges the results in order.
func longestEvenRunParallel(nums []int, parts int) (int, int, int) {
	if len(nums) == 0 || parts < 1 {
		return -1, 0, -1
	}
	parts = min(parts, len(nums))
	size := (len(nums) + parts - 1) / parts
	parts = (len(nums) + size - 1) / size // drop empty trailing chunks

	results := make([]chunkResult, parts)
	var wg sync.WaitGroup
	for p := range parts {
		lo := p * size
		hi := min(lo+size, len(nums))
		wg.Go(func() {
			results[p] = processChunk(nums[lo:hi], lo)
		})
	}
	wg.Wait()

	total := results[0]
	for _, r := range results[1:] {
		total = merge(total, r)
	}

	if total.bestLen <= 2 {
		return -1, 0, -1 // no qualifying subarray found
	}
	return total.bestStart, total.bestLen, total.bestStart + total.bestLen - 1
}

func main() {
	nums := []int{1, 2, 4, 7, 8, 6, 9, 8, 2, 6, 10, 3, 1, 4, 3, 3}

	// 4 parts: [1 2 4 7] [8 6 9 8] [2 6 10 3] [1 4 3 3]
	// The answer 8 2 6 10 crosses the boundary between parts 2 and 3.
	start, length, end := longestEvenRunParallel(nums, 4)

	if length == 0 {
		fmt.Println("no even run longer than 2")
		return
	}

	fmt.Printf("Start: %d, Length: %d, End: %d\n", start, length, end)
	fmt.Println(nums[start : end+1])
}

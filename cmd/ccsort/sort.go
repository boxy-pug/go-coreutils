package main

import (
	"slices"
	"strings"
)

var sortFunctions = map[string]func([]string){
	"stdlib":    stdLibSort,
	"bubble":    bubbleSort,
	"merge":     mergeSort,
	"insertion": insertionSort,
	"quick":     quickSortHelper,
	"selection": selectionSort,
}

// stdLibSort delegates to Go's stdLib sorting implementation.
// The baseline, optimized, stable, handles edge cases
// Time: 0(n log n), space 0(log n)
func stdLibSort(list []string) {
	slices.SortFunc(list, func(a, b string) int {
		return strings.Compare(a, b)
	})
}

// ====> BUBBLE SORT <====

// bubbleSort repeatedly walks through th ewhole list, comparing values next
// to each other, swapping them if they're out of order. That way the largest value
// will "bubble up" to its correct place at the end of unsorted array on each pass.
// Time: 0(n^2), best case 0(n) if sorted with early exit, Space: 0(1), in place.
// Not an efficient algorithm, but easy to understand, good for learning.
func bubbleSort(list []string) {
	swapping := true // flag for early exit, if you do a pass without swapping -> sorted
	end := len(list)

	for swapping {
		swapping = false
		// Walk from start up to the last unsorted value
		for i := 1; i < end; i++ {
			if list[i-1] > list[i] {
				// swap the pair if they're out of order
				list[i-1], list[i] = list[i], list[i-1]
				swapping = true
			}
		}
		// The element at end-1 is now at correct position, so we can decrement end idx
		end--
	}
}

// ====> MERGE SORT <====

// mergeSort is a "divide and conquer" algorithm.
// 1: Divide, split th eslice in two halves.
// 2. "Conquer": recursively sort each half.
// 3. Combine: merge the two sorted halves back into one slice.
// Time: 0(n log n), space 0(n), temp allocations for left and right at each step, but
// it's still proportional to n
func mergeSort(list []string) {
	// The base case for the recursion is: slices of length 0 and 1 are already sorted.
	if len(list) < 2 {
		return
	}

	mid := len(list) / 2
	// Allocate temp slices for two halves.
	// There is a "bottom-up" more advanced merge sort that doesn't allocate
	// at every recursive call, but this top-down version is easier to follow.
	left := make([]string, mid)
	right := make([]string, len(list)-mid)

	// Copy the data over to the slices
	copy(left, list[:mid])
	copy(right, list[mid:])

	// Recursively sort each half
	mergeSort(left)
	mergeSort(right)

	// Merge the sorted halves back into th original list.
	mergeInPlace(list, left, right)
}

// mergeInPlace merges two sorted slices (left and right) into one "list".
// i walks through left, j does right, and k walks through destination list
// len(list) == len(left) + len(right)
// loop continues while both left and right have something in them
// When one side is exhausted the rest of the elements from the other one is copied over.
// (Theyre already sorted so we can copy blindly)
func mergeInPlace(list, left, right []string) {
	i, j, k := 0, 0, 0

	// Merge while both sides have elements in them
	for i < len(left) && j < len(right) {
		// find the val from left or right that is smallest
		// put that into list at current idx, and then increment the
		// idx at the the "winner" left/right and the list idx
		if left[i] <= right[j] {
			list[k] = left[i]
			i++
		} else {
			list[k] = right[j]
			j++
		}
		k++

	}

	// One side is exhausted, copy the rest of the other side
	// Only one of these will execute.
	for i < len(left) {
		list[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		list[k] = right[j]
		j++
		k++
	}
}

// ====> INSERTION SORT <====

// insertionSort builds a sorted portion one element at a time
// everythign before idx i is already sorted. For each new element at i
// we insert it into the correct position in the already sorted portion.
// We do this by swapping it backwards until it is >= the element before it.
// Time: 0(n^2) is worst/average, 0(n) is best case when already sorted. Space is 0(1) in place.
// Metaphor: You pick up one playing card at a time and insert it into the
// correct position among the cards you already hold.
func insertionSort(list []string) {
	for i := range list {
		// j starts at i and walks backwards through the sorted portion
		j := i
		// As long as j > 0 and the preceding element is bigger than current
		// Swap the elements backwards until it finds the right spot
		for j > 0 && list[j-1] > list[j] {
			list[j-1], list[j] = list[j], list[j-1]
			j--
		}
	}
}

// ====> QUICK SORT <====

// quickSortHelper is a thin wrapper for quick sort, exists bcs
// the sorting map stores func([]string), but quick sort needs to pass along
// high and low to recurse properly.
func quickSortHelper(list []string) {
	low := 0
	high := len(list) - 1
	quickSort(list, low, high)
}

// quickSort picks a pivot, partitions the slice so that everything left
// of pivot is smaller and everything to the right is larger, then recursively sort
// thw two partitions. Time: 0(n log n) average, worst case 0(n^2) if bad pivot.
// Not a stable sorting algo, which means it can reorder equal elements across the pivot.
func quickSort(list []string, low, high int) {
	// Base case is low >= high
	if low < high {
		pivot := partition(list, low, high)
		quickSort(list, low, pivot-1)
		quickSort(list, pivot+1, high)
	}
}

// partition rearranges the slice so that all elements < pivot are on the left
// all elements >= pivot are on the right. Picks last element as pivt.
// returns final element of the pivot element
func partition(list []string, low, high int) int {
	pivot := list[high]
	i := low - 1

	for j := low; j < high; j++ {
		if list[j] < pivot {
			i++
			list[i], list[j] = list[j], list[i]
		}
	}
	// Place the pivot between the two regions
	list[i+1], list[high] = list[high], list[i+1]
	return i + 1
}

// ====> SELECTION SORT <====

// selectionSort sort iterates through the slice, finds the minimum element in the unsorted portion,
// and swaps it with the first element of the unsorted portion.
func selectionSort(list []string) {
	// Outer loop: iterates through each element of the slice
	for i := range list {
		// Assume the current index is the minimum
		minIdx := i

		// Inner loop: finds the actual minimum element in the unsorted portion
		for j := i + 1; j < len(list); j++ {
			// Compare the current element with the assumed minimum
			if list[j] < list[minIdx] {
				// Update the minimum index if a smaller element is found
				minIdx = j
			}
		}

		// Swap the found minimum element with the first element of the unsorted portion
		list[i], list[minIdx] = list[minIdx], list[i]
	}
}

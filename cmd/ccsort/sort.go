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

func stdLibSort(list []string) {
	slices.SortFunc(list, func(a, b string) int {
		return strings.Compare(a, b)
	})
}

func bubbleSort(list []string) {
	swapping := true
	end := len(list)

	for swapping {
		swapping = false
		for i := 1; i < end; i++ {
			if list[i-1] > list[i] {
				list[i-1], list[i] = list[i], list[i-1]
				swapping = true
			}
		}
		end--
	}
}

func mergeSort(list []string) {
	if len(list) < 2 {
		return
	}

	mid := len(list) / 2
	left := make([]string, mid)
	right := make([]string, len(list)-mid)

	copy(left, list[:mid])
	copy(right, list[mid:])

	mergeSort(left)
	mergeSort(right)

	mergeInPlace(list, left, right)
}

func mergeInPlace(list, left, right []string) {
	i, j, k := 0, 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			list[k] = left[i]
			i++
		} else {
			list[k] = right[j]
			j++
		}
		k++

	}
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

func insertionSort(list []string) {
	for i := range list {
		j := i
		for j > 0 && list[j-1] > list[j] {
			list[j-1], list[j] = list[j], list[j-1]
			j--
		}
	}
}

func quickSortHelper(list []string) {
	low := 0
	high := len(list) - 1
	quickSort(list, low, high)
}

func quickSort(list []string, low, high int) {
	if low < high {
		pivot := partition(list, low, high)
		quickSort(list, low, pivot-1)
		quickSort(list, pivot+1, high)
	}
}

func partition(list []string, low, high int) int {
	pivot := list[high]
	i := low - 1

	for j := low; j < high; j++ {
		if list[j] < pivot {
			i++
			list[i], list[j] = list[j], list[i]
		}
	}
	list[i+1], list[high] = list[high], list[i+1]
	return i + 1
}

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

package sortalgos

import (
	"slices"
	"strings"

	"github.com/boxy-pug/go-coreutils/cmd/ccsort/config"
)

var SortFunctions = map[config.SortingAlgo]func([]string){
	config.StdLib:    StdLib,
	config.Bubble:    Bubble,
	config.Merge:     Merge,
	config.Insertion: InsertionSort,
	config.Quick:     QuickSortHelper,
	config.Selection: Selection,
}

func StdLib(list []string) {
	slices.SortFunc(list, func(a, b string) int {
		return strings.Compare(a, b)
	})
}

func Bubble(list []string) {
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

func Merge(list []string) {
	if len(list) < 2 {
		return
	}

	mid := len(list) / 2
	left := make([]string, mid)
	right := make([]string, len(list)-mid)

	copy(left, list[:mid])
	copy(right, list[mid:])

	Merge(left)
	Merge(right)

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

func InsertionSort(list []string) {
	for i := range list {
		j := i
		for j > 0 && list[j-1] > list[j] {
			list[j-1], list[j] = list[j], list[j-1]
			j--
		}
	}
}

func QuickSortHelper(list []string) {
	low := 0
	high := len(list) - 1
	QuickSort(list, low, high)
}

func QuickSort(list []string, low, high int) {
	if low < high {
		pivot := partition(list, low, high)
		QuickSort(list, low, pivot-1)
		QuickSort(list, pivot+1, high)
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

// Selection sort iterates through the slice, finds the minimum element in the unsorted portion,
// and swaps it with the first element of the unsorted portion.
func Selection(list []string) {
	// Outer loop: iterates through each element of the slice
	for i := range list {
		// Assume the current index is the minimum
		min_idx := i

		// Inner loop: finds the actual minimum element in the unsorted portion
		for j := i + 1; j < len(list); j++ {
			// Compare the current element with the assumed minimum
			if list[j] < list[min_idx] {
				// Update the minimum index if a smaller element is found
				min_idx = j
			}
		}

		// Swap the found minimum element with the first element of the unsorted portion
		list[i], list[min_idx] = list[min_idx], list[i]
	}
}

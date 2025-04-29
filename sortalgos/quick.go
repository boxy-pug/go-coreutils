package sortalgos

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

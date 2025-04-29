package sortalgos

func InsertionSort(list []string) []string {
	for i := range list {
		j := i
		for j > 0 && list[j-1] > list[j] {
			list[j-1], list[j] = list[j], list[j-1]
			j--
		}
	}
	return list
}

package sortalgos

func Bubble(list []string) []string {
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
		//fmt.Println(list)
		end--
	}

	return list
}

package sortalgos

func MergeSort(list []string) []string {
	if len(list) < 2 {
		return list
	}
	leftHalf := list[:len(list)/2]
	rightHalf := list[len(list)/2:]

	left := MergeSort(leftHalf)
	right := MergeSort(rightHalf)

	return merge(left, right)
}

func merge(left, right []string) []string {
	var final []string
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			final = append(final, left[i])
			i++
		} else {
			final = append(final, right[j])
			j++
		}

	}
	for i < len(left) {
		final = append(final, left[i])
		i++
	}
	for j < len(right) {
		final = append(final, right[j])
		j++
	}
	return final
}

/*
merge_sort() pseudocode

Input: A, an unsorted list of integers

    If the length of A is less than 2, it's already sorted so return it
    Split the input array into two halves down the middle
    Call merge_sort() twice, once on each half
    Return the result of calling merge(sorted_left_side, sorted_right_side) on the results of the merge_sort() calls

merge() pseudocode

Inputs: A and B. Two sorted lists of integers

    Create a new final list of integers.
    Set i and j equal to zero. They will be used to keep track of indexes in the input lists (A and B).
    Use a loop to compare the current elements of A and B. If an element in A is less than or equal to its respective element in B,
	add it to the final list and increment i. Otherwise, add the item in B to the final list and increment j.
	Continue until all items from one of the lists have been added.
    After comparing all the items, there may be some items left over in either A or B. Add those extra items to the final list.
    Return the final list.
*/

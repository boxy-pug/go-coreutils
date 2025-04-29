package sortalgos

import (
	"slices"
	"strings"
)

func StdLib(list []string) []string {
	slices.SortFunc(list, func(a, b string) int {
		return strings.Compare(a, b)
	})
	return list
}

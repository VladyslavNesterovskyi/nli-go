package common

func IntArrayContains(haystack []int, needle int) bool {
	for _, value := range haystack  {
		if needle == value {
			return true
		}
	}
	return false
}

func IntArrayDeduplicate(array []int) []int {
	newArray := []int{}
	previous := map[int]bool{}

	for _, value := range array {
		_, found := previous[value]
		if !found {
			previous[value] = true
			newArray = append(newArray, value)
		}
	}

	return newArray
}
package util

func MaxLength(arr []string) int {
	maxLength := 0
	for _, elem := range arr {
		length := len(elem)
		if maxLength < length {
			maxLength = length
		}
	}
	return maxLength
}

package sequence

func isin(index int, slice []int) bool {
	for _, nextIndex := range slice {
		if index == nextIndex {
			return true
		}
	}
	return false
}

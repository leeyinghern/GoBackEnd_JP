package functions

func CheckIfOverlappingQuestionNumber(a []int, n int) bool {
	// Returns false if there is overlap
	// Returns true otherwise
	for _, val := range a {
		if val == n {
			return false
		}
	}
	return true
}

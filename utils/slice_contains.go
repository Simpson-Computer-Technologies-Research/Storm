package utils

// The SliceContains() function is useds to return
// whether the provided str parameter exists in the
// provided slice paramater
func SliceContains(slice []string, str string) bool {
	// For each value in slice
	for _, v := range slice {
		// If the value equals the provided string
		if v == str {
			return true
		}
	}
	return false
}

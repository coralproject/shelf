package shelf

// stringContains determines if a string value in contained in a slice of strings.
func stringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

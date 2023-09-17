package helpers

// Contains checks if a slice contains an item.
//
// Returns true if present.
func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

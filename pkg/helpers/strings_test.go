package helpers

import (
	"testing"
)

func TestContains(t *testing.T) {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}

	// Test case: element exists in the slice
	exists := "banana"
	result := Contains(slice, exists)
	if !result {
		t.Errorf("Test case failed: expected %s to be found in slice", exists)
	}

	// Test case: element does not exist in the slice
	notExists := "grape"
	result = Contains(slice, notExists)
	if result {
		t.Errorf("Test case failed: unexpected %s found in slice", notExists)
	}
}

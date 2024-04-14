package sketches

import (
	"testing"
)

// Add a lot of items to the sketch
func TestCountMinSketch(t *testing.T) {
	// Initialize a CountMinSketch
	width := 100
	depth := 5
	cms := NewCountMinSketch(uint32(width), uint32(depth))

	// Add some items to the sketch
	items := []string{"apple", "banana", "apple", "cherry", "banana", "cherry", "banana"}
	for _, item := range items {
		cms.Add(item)
	}

	// Test the Count method
	tests := []struct {
		item     string
		expected uint32
	}{
		{"apple", 2},
		{"banana", 3},
		{"cherry", 2},
		{"orange", 0}, // Non-existent item should return 0
	}

	for _, test := range tests {
		count := cms.Count(test.item)
		if count != test.expected {
			t.Errorf("Count of %s returned %d, expected %d", test.item, count, test.expected)
		}
	}
}

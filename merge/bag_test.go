// classify project merge_test.go
package merge

import (
	"testing"
)

func TestContainsIndexes(t *testing.T) {
	b := Bag{
		Content: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}

	if !b.Contains([]int{2, 6, 8}) {
		t.Error("Bag should contain these indexes")
	}

	if b.Contains([]int{20, 60, 80}) {
		t.Error("Bag shouldn't contain these indexes")
	}
}

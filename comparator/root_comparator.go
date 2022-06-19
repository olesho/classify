package comparator

import (
	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
)

type StrictComparator struct {
	a *arena.Arena
}

func NewStrictComparator(a *arena.Arena) *StrictComparator {
	return &StrictComparator{a}
}

func basicCmp(n1, n2 *arena.Node) bool {
	if n1.Id == n2.Id {
		return true
	}
	if n1.Type == n2.Type {
		if n1.Type == html.ElementNode {
			if n1.Data == n2.Data {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

func (c *StrictComparator) Cmp(idx1, idx2 int) float32 {
	if idx1 == 0 || idx2 == 0 {
		return 0
	}

	n1, n2 := c.a.Get(idx1), c.a.Get(idx2)
	for {
		if !basicCmp(n1, n2) {
			return 0
		}
		n1, n2 = c.a.Get(n1.Parent), c.a.Get(n2.Parent)
		if n1.Parent == 0 || n2.Parent == 0 {
			break
		}
	}
	return 1
}

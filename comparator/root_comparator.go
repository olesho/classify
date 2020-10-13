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
	chain1 := c.a.Chain(idx1, 0)
	chain2 := c.a.Chain(idx1, 0)
	size1 := len(chain1)
	size2 := len(chain2)
	if size1 != size2 {
		return 0
	}
	for index := 1; (index < size1) && (index < size2); index++ {
		if !basicCmp(chain1[index], chain2[index]) {
			return 0
		}
	}
	return 1
}

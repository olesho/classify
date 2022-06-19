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

// Cmp calculates similarity between two branches ending with idx1 and idx2
// consumes most RAM; needs optimization
//func (c *StrictComparator) CmpUnoptimized(idx1, idx2 int) float32 {
//	chain1 := c.a.Chain(idx1, 0)
//	chain2 := c.a.Chain(idx2, 0)
//	size1 := len(chain1)
//	size2 := len(chain2)
//	if size1 != size2 {
//		return 0
//	}
//	for index := 0; (index < size1) && (index < size2); index++ {
//		if !basicCmp(chain1[index], chain2[index]) {
//			return 0
//		}
//	}
//	return 1
//}

func (c *StrictComparator) Cmp(idx1, idx2 int) float32 {
	//if (idx1 == 7 || idx1 == 22 || idx1 == 37 || idx1 == 51 || idx1 == 66 || idx1 == 81 || idx1 == 95) &&
	//	(idx2 == 7 || idx2 == 22 || idx2 == 37 || idx2 == 51 || idx2 == 66 || idx2 == 81 || idx2 == 95) {
	//	fmt.Println()
	//}

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

	//c := []int{nId}
	//for n := a.Get(nId); n.Parent != 0; n = a.Get(n.Parent) {
	//	c = append(c, n.Parent)
	//}
	//return c

	//chain1 := c.a.ChainIDXs(idx1, 0)
	//chain2 := c.a.ChainIDXs(idx2, 0)

	//size1 := len(chain1)
	//size2 := len(chain2)
	//if size1 != size2 {
	//	return 0
	//}
	//for index := 0; (index < size1) && (index < size2); index++ {
	//	if !basicCmp(c.a.Get(chain1[index]), c.a.Get(chain2[index])) {
	//		return 0
	//	}
	//}
	//return 1
}

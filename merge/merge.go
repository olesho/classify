// classify project merge.go
package merge

import (
	"fmt"
	"sort"

	"github.com/olesho/classify"
)

type MergeClassificator struct {
	*classify.Arena
	bags Bags
}

func NewMergeClassificator(a *classify.Arena) MergeClassificator {
	bags := make([]Bag, len(a.List))
	for i, _ := range a.List {
		bags[i] = Bag{
			Content: []int{i},
			//Arena:   a,
		}
	}
	return MergeClassificator{
		Arena: a,
		bags:  Bags{bags},
	}
}

func (c *MergeClassificator) Bags() Bags {
	return c.bags
}

func (c *MergeClassificator) cmp(n1, n2 int) float64 {
	bag1 := c.bags.List[n1]
	bag2 := c.bags.List[n2]

	if len(bag1.Content) == 0 {
		return 0
	}

	if len(bag2.Content) == 0 {
		return 0
	}

	arena1 := bag1.Arena
	index1 := 0
	if arena1 == nil {
		arena1 = c.Arena
		index1 = n1
	}
	arena2 := bag2.Arena
	index2 := 0
	if arena2 == nil {
		arena2 = c.Arena
		index2 = n2
	}

	r := classify.CmpDeepRate(arena1, arena2, index1, index2)
	if r != nil {
		return r.Rate()
	}
	return 0

}

func (c *MergeClassificator) merge(n1, n2 int) bool {
	bag1 := c.bags.List[n1]
	bag2 := c.bags.List[n2]

	arena1 := c.bags.List[n1].Arena
	index1 := 0
	if arena1 == nil {
		arena1 = c.Arena
		index1 = n1
	}

	arena2 := c.bags.List[n2].Arena
	index2 := 0
	if arena2 == nil {
		arena2 = c.Arena
		index2 = n2
	}

	//fmt.Println("Merge bags:", n1, "<=", n2)

	newArena := classify.Merge(arena1, arena2, index1, index2)
	if len(newArena.List) > 0 {
		newBag := Bag{
			Arena:   newArena,
			Content: append(bag1.Content, bag2.Content...),
		}

		if newBag.Efficacy() > c.bags.List[n1].Efficacy() {
			c.bags.List[n1] = newBag
			c.bags.List[n2].Clear()

			return true
		}
	}

	return false
}

/*
func (a *MergeClassificator) bagNested(nestedBag, inBag Bag) bool {
	cnt := 0

	if len(nestedBag.Content) != len(inBag.Content) {
		return false
	}

	for _, nNested := range nestedBag.Content {
		if !a.nodeNested(nNested, inBag) {
			return false
		}
		cnt++
	}
	if cnt == len(nestedBag.Content) && nestedBag.Rate < inBag.Rate {
		return true
	}
	return false
}


func (a *MergeClassificator) nodeNested(nestedNode int, inBag Bag) bool {
	for _, bn := range inBag.Content {
		if a.pathNested(nestedNode, bn) {
			return true
		}
	}

	return false
}


func (a *MergeClassificator) pathNested(inNode, nestedNode int) bool {
	path := a.PathArray(inNode)
	for _, item := range path {
		if item == nestedNode {
			return true
		}
	}
	return false
}

func (c *MergeClassificator) filterNested() {
	for i1, b1 := range c.bags.List {
		for i2, b2 := range c.bags.List {
			if i1 != i2 {
				if c.bagNested(b1, b2) {
					c.bags.List[i1].Clear()
				}
			}
		}
	}
}

*/
func (c *MergeClassificator) Run() {
	for n1, _ := range c.List {
		bestBagIndex, maxRate := c.findBestFit(n1, n1+1)

		if n1 == 134 {
			fmt.Println("<<<", bestBagIndex, "Rate:", maxRate)
			spsd := c.cmp(n1, 400)
			fmt.Println("Supposed rate", spsd)
		}

		for bestBagIndex > -1 && maxRate > 0 {
			if !c.merge(n1, bestBagIndex) {
				break
			}
			bestBagIndex, maxRate = c.findBestFit(n1, n1+1)
		}

		if len(c.bags.List[n1].Content) < 2 {
			c.bags.List[n1].Clear()
		}
	}

	//c.filterNested()
	sort.Sort(c.bags)
}

func (c *MergeClassificator) findBestFit(n1 int, offset int) (bestBagIndex int, maxRate float64) {
	bestBagIndex = -1
	for n2 := offset; n2 < len(c.List); n2++ {
		if n2 != n1 {
			nextRate := c.cmp(n1, n2)
			if nextRate > maxRate {
				maxRate = nextRate
				bestBagIndex = n2
			}
		}
	}
	return
}

func (c *MergeClassificator) BagsContaining(indexes []int) []Bag {
	res := make([]Bag, 0)
	for _, b := range c.bags.List {
		if b.Contains(indexes) {
			res = append(res, b)
		}
	}
	return res
}

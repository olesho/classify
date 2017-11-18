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

/*
func (c *MergeClassificator) Classify(n int) {
	var maxRate float64 = 0
	var bestBagIndex = -1
	for i, bag := range c.bags.List {
		if len(bag.Content) > 0 {
			if i != n {
				var r *classify.CmpResult = nil
				if bag.Arena == nil {
					r = classify.CmpDeepRate(c.Arena, c.Arena, i, n)
				} else {
					r = classify.CmpDeepRate(bag.Arena, c.Arena, 0, n)
				}
				if r != nil {
					val := r.Rate()
					if val > maxRate {
						maxRate = val
						bestBagIndex = i
					}
				}
			}
		}
	}
	c.put(n, bestBagIndex, maxRate)
}
*/

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

	fmt.Println("Merge bags:", n1, "<=", n2)
	newArena := classify.Merge(arena1, arena2, index1, index2)
	if len(newArena.List) > 0 {
		newBag := Bag{
			Arena:   newArena,
			Content: append(bag1.Content, bag2.Content...),
		}

		if newBag.Efficacy() > c.bags.List[n1].Efficacy() {
			fmt.Println(newBag.Efficacy())
			c.bags.List[n1] = newBag
			c.bags.List[n2].Clear()
			return true
		}
	}

	return false

	/*
			if c.bags.List[bestBagIndex].Arena == nil {
				fmt.Println("Merge nodes:", n, "<=>", bestBagIndex)
				newArena = classify.Merge(c.Arena, c.Arena, bestBagIndex, n)
				content = []int{n, bestBagIndex}
			} else {
				fmt.Println("Merge bag:", bestBagIndex, "<=>", n)
				newArena = classify.Merge(c.bags.List[bestBagIndex].Arena, c.Arena, 0, n)
				content = append(c.bags.List[bestBagIndex].Content, n)
			}


		if len(newArena.List) > 0 {
			maxResult = Bag{
				Arena:   newArena,
				Content: content,
			}
			if maxResult.Efficacy() > c.bags.List[bestBagIndex].Efficacy() {
				fmt.Println(maxResult.Efficacy())
				c.bags.List[bestBagIndex] = maxResult
				c.bags.List[n].Clear()
			}

			return
		}
	*/
}

/*
// puts new comparation result for a node into bag
func (c *MergeClassificator) put(n int, bestBagIndex int, maxRate float64) {
	// try to put into existing bag
	if bestBagIndex > -1 && maxRate > 0 {
		var maxResult Bag
		var newArena *classify.Arena
		var content []int

		if c.bags.List[bestBagIndex].Arena == nil {
			fmt.Println("Merge nodes:", n, "<=>", bestBagIndex)
			newArena = classify.Merge(c.Arena, c.Arena, bestBagIndex, n)
			content = []int{n, bestBagIndex}
		} else {
			fmt.Println("Merge bag:", bestBagIndex, "<=>", n)
			newArena = classify.Merge(c.bags.List[bestBagIndex].Arena, c.Arena, 0, n)
			content = append(c.bags.List[bestBagIndex].Content, n)
		}

		if len(newArena.List) > 0 {
			maxResult = Bag{
				Arena:   newArena,
				Content: content,
			}
			if maxResult.Efficacy() > c.bags.List[bestBagIndex].Efficacy() {
				fmt.Println(maxResult.Efficacy())
				c.bags.List[bestBagIndex] = maxResult
				c.bags.List[n].Clear()
			}

			return
		}
	}
	c.bags.List[n].Clear()
}
*/

func (a *MergeClassificator) bagNested(nestedBag, inBag classify.Bag) bool {
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

func (a *MergeClassificator) nodeNested(nestedNode int, inBag classify.Bag) bool {
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

func (c *MergeClassificator) Run() {
	for n1, _ := range c.List {
		bestBagIndex, maxRate := c.findBestFit(n1, n1+1)
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

/*
func (c *MergeClassificator) Run() {
	for n, _ := range c.List {
		if len(c.bags.List[n].Content) > 0 {
			c.Classify(n)
		}
	}

	for i := len(c.List); i < c.bags.Len(); i++ {
		if len(c.bags.List[n].Content) > 0 {
			c.Classify(n)
		}
	}

	sort.Sort(c.bags)
}
*/

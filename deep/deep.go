// classify project deep.go
package deep

import (
	"sort"

	"github.com/olesho/classify"
)

type DeepClassificator struct {
	classify.Arena
	bags classify.Bags
}

func NewDeepClassificator(a *classify.Arena) DeepClassificator {
	bags := make([]classify.Bag, len(a.List))
	for i, _ := range a.List {
		bags[i].Content = []int{i}
	}
	return DeepClassificator{
		Arena: *a,
		bags:  classify.Bags{bags},
	}
}

func (c *DeepClassificator) Bags() classify.Bags {
	return c.bags
}

func (c *DeepClassificator) Classify(n int) {
	var maxRate float64 = 0
	var bestBagIndex = -1
	var maxResult *classify.CmpResult
	for i, bag := range c.bags.List {
		if len(bag.Content) > 0 {
			if bag.Content[0] != n {
				r := c.CmpDeepRate(bag.Content[0], n)
				if r != nil {
					val := r.Rate()
					if val > maxRate {
						maxRate = val
						bestBagIndex = i
						maxResult = r
					}
				}
			}
		}
	}
	//if maxResult != nil {
	c.put(n, bestBagIndex, maxRate, maxResult)
	//}
}

// puts new comparation result for a node into bag
func (c *DeepClassificator) put(n int, bestBagIndex int, maxRate float64, maxResult *classify.CmpResult) {
	// try to put into existing bag
	if bestBagIndex > -1 && maxResult != nil {
		if maxRate*float64(len(c.bags.List[bestBagIndex].Content)+1) > c.bags.List[bestBagIndex].Rate {
			c.bags.List[bestBagIndex].Content = append(c.bags.List[bestBagIndex].Content, n)
			c.bags.List[bestBagIndex].Rate = maxRate * float64(len(c.bags.List[bestBagIndex].Content))

			if len(c.bags.List[bestBagIndex].Content) == 2 {
				// initial Sum
				c.bags.List[bestBagIndex].Sum = maxResult.Sum * 2
				c.bags.List = append(c.bags.List, c.bags.List[bestBagIndex])

				c.bags.List[bestBagIndex].Clear()
				c.bags.List[n].Clear()

				return
			}
			c.bags.List[bestBagIndex].Sum = c.bags.List[bestBagIndex].Sum + maxResult.Sum
			return
		}
	}

	// create new bag
	/*
		newBag := classify.Bag{
			Content: []int{n},
			Rate:    maxRate,
		}

		// get rid of nested bags

			for i, b := range c.bags.List {
				// remove nested bag
				if c.bagNested(newBag, b) {
					return
				}

				// or replace nested bag
				if c.bagNested(b, newBag) {
					c.bags.List[i] = newBag
					return
				}
			}


		// append new bag
		c.bags.List = append(c.bags.List, newBag)
	*/
}

func (c *DeepClassificator) filterNested() {
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

func (a *DeepClassificator) bagNested(nestedBag, inBag classify.Bag) bool {
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

func (a *DeepClassificator) nodeNested(nestedNode int, inBag classify.Bag) bool {
	for _, bn := range inBag.Content {
		if a.pathNested(nestedNode, bn) {
			return true
		}
	}

	return false
}

func (a *DeepClassificator) pathNested(inNode, nestedNode int) bool {
	path := a.PathArray(inNode)
	for _, item := range path {
		if item == nestedNode {
			return true
		}
	}
	return false
}

func (c *DeepClassificator) Run() {
	for n, _ := range c.List {
		if len(c.bags.List[n].Content) > 0 {
			c.Classify(n)
		}
	}
	c.filterNested()
	sort.Sort(c.bags)
}

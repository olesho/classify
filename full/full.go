// classify project full.go
package full

import (
	"fmt"
	//	"sort"

	"github.com/olesho/classify"
)

type FullClassificator struct {
	classify.Arena
	bags Bags
}

type Bags struct {
	List []Bag
}

type Bag struct {
	*classify.Arena
	Content []int
}

func NewFullClassificator(a *classify.Arena) FullClassificator {
	bags := make([]Bag, len(a.List))
	for i, _ := range a.List {
		//bags[i].Arena = a
		bags[i].Content = []int{i}
	}
	return FullClassificator{
		//Arena: *a,
		bags: Bags{bags},
	}
}

func (c *FullClassificator) Classify(n int) {
	var maxRate float64 = 0
	var bestBagIndex = -1
	var maxResult *classify.CmpResult
	for i, bag := range c.bags.List {
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

	fmt.Println(bestBagIndex, maxResult)

	/*
		if bestBagIndex > -1 && maxResult != nil {
			newArena := c.Arena.Clone(n)
			c.bags.List[bestBagIndex].Arena = classify.Merge(newArena, c.bags.List[bestBagIndex].Arena, 0, n)
			c.bags.List[bestBagIndex].Content = append(c.bags.List[bestBagIndex].Content, n)
		}
	*/
}

func (c *FullClassificator) Run() {
	c.Classify(0)
	c.Classify(1)
	c.Classify(2)
	/*
		for n, _ := range c.List {
			fmt.Println(n)
			c.Classify(n)
		}
	*/
}

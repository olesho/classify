// classify project simple.go
package simple

import (
	"github.com/olesho/classify"
)

type SimpleClassificator struct {
	classify.Arena
	bags classify.Bags
}

func NewSimpleClassificator(a *classify.Arena) SimpleClassificator {
	return SimpleClassificator{
		Arena: *a,
		bags:  classify.Bags{[]classify.Bag{}},
	}
}

func (c *SimpleClassificator) Bags() classify.Bags {
	return c.bags
}

func (c *SimpleClassificator) Classify(n int) {
	var max float64 = 0
	var max_i = -1
	for i, bag := range c.bags.List {
		res := c.CmpColumn(bag.Content[0], n)
		if res != nil {
			val := res.Rate()
			if val > max {
				max = val
				max_i = i
			}
		}
	}

	if max_i > -1 {
		c.bags.List[max_i].Content = append(c.bags.List[max_i].Content, n)
	} else {
		c.bags.List = append(c.bags.List, classify.Bag{
			Content: []int{n},
		})
	}
}

func (c *SimpleClassificator) Run() {
	for n, _ := range c.List {
		c.Classify(n)
	}

	for i := 0; i < c.bags.Len()-1; i++ {
		if classify.IsParent(&c.Arena, c.bags.List[i], c.bags.List[i+1]) {
			c.bags.List = append(c.bags.List[:i+1], c.bags.List[i+2:]...)
			i--
		}
	}
}

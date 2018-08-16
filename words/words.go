// classify project simple.go
package words

import (
	"github.com/olesho/classify"
)

type WordsClassificator struct {
	classify.Arena
	bags classify.Bags
}

func NewWordsClassificator(a *classify.Arena) WordsClassificator {
	return WordsClassificator{
		Arena: *a,
		bags:  classify.Bags{[]classify.Bag{}},
	}
}

func (c *WordsClassificator) Bags() classify.Bags {
	return c.bags
}

func (c *WordsClassificator) Classify(n int) {
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

func (c *WordsClassificator) Run() {
	for n, _ := range c.List {
		c.Classify(n)
	}
}

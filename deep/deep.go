// classify project deep.go
package deep

import (
	//	"fmt"
	"sort"

	"github.com/olesho/classify"
)

type DeepClassificator struct {
	classify.Arena
	bags classify.Bags
}

func NewDeepClassificator(a *classify.Arena) DeepClassificator {
	return DeepClassificator{
		Arena: *a,
		bags:  classify.Bags{[]classify.Bag{}},
	}
}

func (c *DeepClassificator) Bags() classify.Bags {
	return c.bags
}

func (c *DeepClassificator) Classify(n int) {
	var max float64 = 0
	var max_i = -1
	var max_r *classify.CmpResult
	for i, bag := range c.bags.List {
		r := c.CmpDeepRate(bag.Content[0], n)
		if r != nil {
			val := r.Result()
			if val > max {
				max = val
				max_i = i
				max_r = r
			}
		}
	}
	//if max_r != nil {
	c.put(n, max_i, max, max_r)
	//}
}

func (c *DeepClassificator) put(n int, max_i int, max float64, max_r *classify.CmpResult) {
	if max_i > -1 && max_r != nil {
		if max*float64(len(c.bags.List[max_i].Content)+1) > c.bags.List[max_i].Rate {
			c.bags.List[max_i].Content = append(c.bags.List[max_i].Content, n)
			c.bags.List[max_i].Rate = max * float64(len(c.bags.List[max_i].Content))
			c.bags.List[max_i].Efficacy = c.bags.List[max_i].Efficacy + max_r.Rate
			return
		}
	}

	c.bags.List = append(c.bags.List, classify.Bag{
		Content:  []int{n},
		Efficacy: 0, //* max_r.Count,
		Rate:     max,
	})

}

func (c *DeepClassificator) Run() {
	for n, _ := range c.List {
		c.Classify(n)
	}
	sort.Sort(c.bags)
}

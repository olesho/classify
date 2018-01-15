// classify project merge.go
package advanced

import (
	"github.com/olesho/classify"
)

type Bag struct {
	Content []int
	Sum     int
	Rate    float64
	*classify.Arena
}

func (b *Bag) Clear() {
	b.Content = []int{}
}

func (b *Bag) Efficacy() int {
	if b.Arena == nil {
		return 0
	}
	if len(b.Arena.List) > 1 {
		return b.Arena.Rate(0) * len(b.Content)
	}
	return 0
}

func (b *Bag) Contains(indexes []int) bool {
	for _, index := range indexes {
		if !b.ContainsIndex(index) {
			return false
		}
	}
	return true
}

func (b *Bag) ContainsIndex(index int) bool {
	for _, n := range b.Content {
		if n == index {
			return true
		}
	}
	return false
}

type Bags struct {
	List []Bag
}

func (b Bags) Len() int {
	return len(b.List)
}

func (b Bags) Less(i, j int) bool {
	//return len(b.List[i].Content) > len(b.List[j].Content)
	return b.List[i].Efficacy() > b.List[j].Efficacy()
}

func (b Bags) Swap(i, j int) {
	b.List[i], b.List[j] = b.List[j], b.List[i]
}

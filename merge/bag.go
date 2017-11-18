// classify project merge.go
package merge

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
		return b.Arena.Rate(1) * len(b.Content)
	}
	return 0
}

type Bags struct {
	List []Bag
}

func (b Bags) Len() int {
	return len(b.List)
}

func (b Bags) Less(i, j int) bool {
	return b.List[i].Efficacy() > b.List[j].Efficacy()
}

func (b Bags) Swap(i, j int) {
	b.List[i], b.List[j] = b.List[j], b.List[i]
}

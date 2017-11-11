// classify project classify.go
package classify

type Bag struct {
	Content  []int
	Efficacy int
	Rate     float64
}

type Bags struct {
	List []Bag
}

func (b Bags) Len() int {
	return len(b.List)
}

func (b Bags) Less(i, j int) bool {
	return b.List[i].Efficacy > b.List[j].Efficacy
}

func (b Bags) Swap(i, j int) {
	b.List[i], b.List[j] = b.List[j], b.List[i]
}

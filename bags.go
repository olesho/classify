// classify project classify.go
package classify

type Bag struct {
	Content []int
	Sum     int
	Rate    float64
}

func (b *Bag) Clear() {
	b.Content = []int{}
	b.Rate = 0
	b.Sum = 0
}

func (b Bag) Efficacy() float64 {
	//return float64(b.Sum) / float64(len(b.Content))
	return float64(b.Sum)
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

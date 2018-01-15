// classify project advanced.go
package advanced

import (
	//"fmt"
	"sort"

	"github.com/olesho/classify"
)

type RoundMatrix [][]interface{}

func NewRoundMatrix(n int) RoundMatrix {
	m := make([][]interface{}, n)
	for i := 0; i < n; i++ {
		m[i] = make([]interface{}, n-i-1)
	}
	return m
}

func (rm RoundMatrix) Get(x, y int) interface{} {
	if x == y {
		panic("Round Matrix Comparison: X cannot equal Y")
	}
	if x < y {
		return [][]interface{}(rm)[x][y-x-1]
	}
	return [][]interface{}(rm)[y][x-y-1]
}

func (rm RoundMatrix) Set(x, y int, val interface{}) {
	if x == y {
		panic("Round Matrix Comparison: X cannot equal Y")
	}
	if x < y {
		[][]interface{}(rm)[x][y-x-1] = val
		return
	}
	[][]interface{}(rm)[y][x-y-1] = val
	return
}

/*
type Item struct {
	x   int
	y   int
	val float64
}
*/

type AdvancedClassificator struct {
	// arena containing all nodes
	*classify.Arena
	// list of all bags (initially equals to number of all nodes)
	bags Bags
	// matrix of comparation of all nodes
	matrix RoundMatrix
	// sorted matrix
	//sorted []Item
}

func NewAdvancedClassificator(a *classify.Arena) AdvancedClassificator {
	bags := make([]Bag, len(a.List))
	for i, _ := range a.List {
		bags[i] = Bag{
			Content: []int{i},
			//Arena:   a,
		}
	}
	return AdvancedClassificator{
		Arena: a,
		bags:  Bags{bags},
	}
}

func (c *AdvancedClassificator) Bags() Bags {
	return c.bags
}

func (c *AdvancedClassificator) cmp(n1, n2 int) int {
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
		return r.Sum
		//return r.Rate()
	}
	return 0

}

func (c *AdvancedClassificator) merge(n1, n2 int) bool {
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

	//fmt.Println("Merge bags:", n1, "<=", n2)

	newArena := classify.Merge(arena1, arena2, index1, index2)
	if len(newArena.List) > 0 {
		newBag := Bag{
			Arena:   newArena,
			Content: append(bag1.Content, bag2.Content...),
		}

		if newBag.Efficacy() > c.bags.List[n1].Efficacy() {
			if !c.nested(newBag) {
				c.bags.List[n1] = newBag
			} else {
				for i, b := range c.bags.List {
					if c.bagNested(b, newBag) {
						c.bags.List[i] = newBag
						break
					}
				}
				c.bags.List[n1].Clear()
			}

			c.bags.List[n2].Clear()

			for n, _ := range c.List {
				if n != n2 {
					c.matrix.Set(n, n2, int(0))
				}
			}

			c.findAll(n1)

			return true
		}
	}

	c.matrix.Set(n1, n2, int(0))

	return false
}

func (c *AdvancedClassificator) Run() {
	c.matrix = NewRoundMatrix(len(c.List))

	for n1, _ := range c.List {
		c.findRow(n1, n1+1)
	}

	x, y, v := c.findBest()
	for v > 200 {
		c.merge(x, y)
		x, y, v = c.findBest()
	}

	//	c.filterNested()
	sort.Sort(c.bags)
}

// compare all nodes (since offset) to n1 and write results to matrix
func (c *AdvancedClassificator) findRow(n1 int, offset int) {
	for n2 := offset; n2 < len(c.List); n2++ {
		if n2 != n1 {
			nextRate := c.cmp(n1, n2)
			c.matrix.Set(n1, n2, nextRate)
		}
	}
	return
}

// compare all nodes to n1 and write results to matrix
func (c *AdvancedClassificator) findAll(n1 int) {
	for n2, _ := range c.List {
		if n2 != n1 {
			nextRate := c.cmp(n1, n2)
			c.matrix.Set(n1, n2, nextRate)
		}
	}
	return
}

func (c *AdvancedClassificator) findBest() (best_x, best_y, val int) {
	for x, row := range [][]interface{}(c.matrix) {
		for y, v := range row {
			v := v.(int)
			if v > val {
				best_x = x
				best_y = y + x + 1
				val = v
			}
		}
	}
	return
}

func (c *AdvancedClassificator) BagsContaining(indexes []int) []Bag {
	res := make([]Bag, 0)
	for _, b := range c.bags.List {
		if b.Contains(indexes) {
			res = append(res, b)
		}
	}
	return res
}

func (a *AdvancedClassificator) bagNested(nestedBag, inBag Bag) bool {
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

func (a *AdvancedClassificator) nodeNested(nestedNode int, inBag Bag) bool {
	for _, bn := range inBag.Content {
		if a.pathNested(nestedNode, bn) {
			return true
		}
	}

	return false
}

func (a *AdvancedClassificator) pathNested(inNode, nestedNode int) bool {
	path := a.PathArray(inNode)
	for _, item := range path {
		if item == nestedNode {
			return true
		}
	}
	return false
}

func (c *AdvancedClassificator) nested(bag Bag) bool {
	for _, inBag := range c.bags.List {
		if c.bagNested(bag, inBag) {
			return true
		}
	}
	return false
}

/*
func (c *AdvancedClassificator) filterNested() {
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
*/

package bags

import (
	"sort"
	"strings"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

type Processor struct {
	threshold  float64
	bags       []*Bag
	comparator func(*classify.Arena, *classify.Node, *classify.Node) float64
	arena      *classify.Arena
}

func New(a *classify.Arena, threshold float64, comparator func(*classify.Arena, *classify.Node, *classify.Node) float64) *Processor {
	return &Processor{
		threshold:  threshold,
		comparator: comparator,
		arena:      a,
	}
}

func (p *Processor) compareToBag(bag *Bag, c *classify.Node) float64 {
	total := .0
	for _, item := range bag.Nodes {
		total += p.comparator(p.arena, item, c)
	}
	return total / float64(len(bag.Nodes))
}

func (p *Processor) Next(c *classify.Node) {
	max := .0
	max_i := 1
	for i, bag := range p.bags {
		val := p.compareToBag(bag, c)

		if val > max {
			max = val
			max_i = i
		}
	}

	if (max_i > -1) && (max >= p.threshold) {
		p.bags[max_i].Nodes = append(p.bags[max_i].Nodes, c)
	} else {
		p.bags = append(p.bags, &Bag{[]*classify.Node{c}, 0})
	}
}

func StrictComparator(a *classify.Arena, n1, n2 *classify.Node) float64 {
	if n1.Type == n2.Type {
		if n2.Type == html.ElementNode && n1.Data != n2.Data {
			return 0
		}
		return 1
	}
	return 0
}

//const textPoints int = 5
//const nodePoints int = 5
func elementComparator(a *classify.Arena, n1, n2 *classify.Node) float64 {
	if StrictComparator(a, n1, n2) == 0 {
		return 0
	}
	return cmpAttr(n1.Attr, n2.Attr)
}

func ColumnComparator(a *classify.Arena, n1, n2 *classify.Node) float64 {
	sum := elementComparator(a, n1, n2) * 0.2
	if sum == 0 {
		return 0
	}

	chain1 := a.Chain(n1.Id, 0)
	chain2 := a.Chain(n2.Id, 0)

	size1 := len(chain1)
	size2 := len(chain2)
	for size1 > 0 && size2 > 0 {
		size1--
		size2--

		next_rate := elementComparator(a, chain1[size1], chain2[size2])
		if next_rate == 0 {
			break
		}
		sum += next_rate
	}

	return (sum * 2) / float64(len(chain1)+len(chain2)+2)
}

func ExtendedComparator(a *classify.Arena, n1, n2 *classify.Node) float64 {
	weightedRate := ColumnComparator(a, n1, n2) * 0.1
	if weightedRate == 0 {
		return 0
	}

	return weightedRate + ChildComparator(a, n1, n2)*0.9
}

func sum(vals []float64) (r float64) {
	for _, v := range vals {
		r += v
	}
	return
}

func ChildComparator(a *classify.Arena, n1, n2 *classify.Node) float64 {
	size1, size2 := len(n1.Children), len(n2.Children)
	rates1, rates2 := make([]float64, size1), make([]float64, size2)
	for i1, idx1 := range n1.Children {
		for i2, idx2 := range n2.Children {
			rate := elementComparator(a, a.Get(idx1), a.Get(idx2))
			if rate > rates1[i1] {
				rates1[i1] = rate
			}
			if rate > rates2[i2] {
				rates2[i2] = rate
			}
		}
	}

	return (sum(rates1) + sum(rates2)) / float64(size1+size2)
}

const attrKeyPoints = .5
const attrValPoints = .5

func cmpAttr(attr1, attr2 []html.Attribute) float64 {
	if len(attr1) == 0 && len(attr2) == 0 {
		return 1
	}

	attrSum := .0
	totalAttr := .0
	for _, a1 := range attr1 {
		if a1.Key != "class" {
			totalAttr += 1
		}
	}

	for _, a2 := range attr2 {
		if a2.Key != "class" {
			totalAttr += 1
		}
	}

	for _, a1 := range attr1 {
		for _, a2 := range attr2 {
			if a1.Key == a2.Key && a1.Key != "class" && a2.Key != "class" {
				attrSum += attrKeyPoints
				attrSum += cmpStrings(a1.Val, a2.Val) * attrValPoints
			}
		}
	}

	classesSum := .0
	classes1 := []string{}
	classes2 := []string{}
	for _, a1 := range attr1 {
		if a1.Key == "class" {
			classes1 = strings.Fields(a1.Val)
		}
	}
	for _, a2 := range attr2 {
		if a2.Key == "class" {
			classes2 = strings.Fields(a2.Val)
		}
	}

	for _, c1 := range classes1 {
		for _, c2 := range classes2 {
			if c1 == c2 {
				classesSum += 1
			}
		}
	}

	return (float64(classesSum*2) + float64(attrSum*2)) / (float64(len(classes1)+len(classes2)) + totalAttr)
}

func cmpStrings(s1 string, s2 string) float64 {
	if len(s1) == 0 && len(s2) == 0 {
		return 1
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	var r int
	l := len(s2)
	if len(s1) < len(s2) {
		l = len(s1)
	}

	for i := 0; i < l; i++ {
		if s1[i] == s2[i] {
			r++
		} else {
			break
		}
	}
	return float64(r*2) / float64(len(s1)+len(s2))
}

func CalcVol(a *classify.Arena, n *classify.Node) int {
	vol := 0
	for _, childIdx := range n.Children {
		vol += CalcVol(a, a.Get(childIdx))
		if n.Type == html.TextNode {
			vol += len(n.Data)
		} else {
			vol += 1
		}
	}
	n.Volume = vol
	return vol
}

type Bag struct {
	Nodes  []*classify.Node
	Volume int
}

type BagGroup struct {
	Bags   []*Bag
	Volume int
}

func belongs(a *classify.Arena, nodes, toNodes []*classify.Node) bool {
	for i := range nodes {
		if !a.HasParent(nodes[i].Id, toNodes[i].Id) {
			return false
		}
	}
	return true
}

func intersects(a *classify.Arena, ns1, before []*classify.Node) bool {
	for i := range ns1[:len(ns1)-1] {
		// ns1[i].Id is between ns2[i].Id and ns2[i+1].Id
		if !(ns1[i].Id > before[i].Id && ns1[i].Id < before[i+1].Id) {
			return false
		}
	}
	if !(ns1[len(ns1)-1].Id > before[len(ns1)-1].Id) {
		return false
	}
	return true
}

func intersectsLosely(a *classify.Arena, ns1, before []*classify.Node) bool {

	for i := range ns1[:len(ns1)-1] {
		// ns1[i].Id is between ns2[i].Id and ns2[i+1].Id
		if !(ns1[i].Id > before[i].Id && ns1[i].Id < before[i+1].Id) {
			return false
		}
	}
	return true
}
func groupBags(a *classify.Arena, bags []*Bag) []*BagGroup {
	for i1, bag1 := range bags {
		if i1+1 < len(bags) {
			for i2, bag2 := range bags[i1+1:] {
				if len(bag1.Nodes) == len(bag2.Nodes) {
					if belongs(a, bag1.Nodes, bag2.Nodes) {
						// remove bag2
						bags = append(bags[:i2], bags[i2+1:]...)
					} else if belongs(a, bag2.Nodes, bag1.Nodes) {
						// remove bag1
						bags = append(bags[:i1], bags[i1+1:]...)
					}
				}
			}
		}
	}

	groups := []*BagGroup{}
	for _, bag1 := range bags {
		groups = checkNextIntersectionStrict(a, groups, bag1)
	}

	return groups
}

func checkNextIntersectionStrict(a *classify.Arena, groups []*BagGroup, bag1 *Bag) []*BagGroup {
	for _, g := range groups {
		if len(bag1.Nodes) == len(g.Bags[0].Nodes) {
			for i2, bag2 := range g.Bags {
				if intersects(a, bag2.Nodes, bag1.Nodes) {
					// add bag1
					g.Bags = append(g.Bags[:i2], append([]*Bag{bag1}, g.Bags[i2:]...)...)
					g.Volume += bag1.Volume
					return groups
				} else if intersects(a, bag1.Nodes, bag2.Nodes) {
					// add bag1
					g.Bags = append(g.Bags, bag1)
					g.Volume += bag1.Volume
					return groups
				}
			}
		}
	}
	groups = append(groups, &BagGroup{Bags: []*Bag{bag1}, Volume: bag1.Volume})
	return groups
}

/*
// add that later
func checkNextIntersectionLose(a *classify.Arena, groups []*BagGroup, bag1 *Bag) []*BagGroup {
	for _, g := range groups {
		if len(bag1.Nodes) < len(g.Bags[0].Nodes) {
			for i2, bag2 := range g.Bags {
				if intersects(a, bag2.Nodes, bag1.Nodes) {
					g.Bags = append(g.Bags[:i2], append([]*Bag{bag1}, g.Bags[i2:]...)...)
					return groups
				} else if intersects(a, bag1.Nodes, bag2.Nodes) {
					g.Bags = append(g.Bags, bag1)
					return groups
				}
			}
		}
	}
	groups = append(groups, &BagGroup{Bags: []*Bag{bag1}})
	return groups
}
*/

func transpose(group *BagGroup) [][]*classify.Node {
	size := len(group.Bags[0].Nodes)
	newGroup := make([][]*classify.Node, size)
	for i := 0; i < size; i++ {
		row := []*classify.Node{}
		for _, bag := range group.Bags {
			row = append(row, bag.Nodes[i])
		}
		newGroup[i] = row
	}
	return newGroup
}

func Parse(arena *classify.Arena) ([][][]*classify.Node, error) {
	CalcVol(arena, arena.Get(0))

	// group by tags && types coincided
	p := New(arena, 1, StrictComparator)
	for _, n := range arena.List {
		p.Next(n)
	}

	finalBags := []*Bag{}

	for idx, bag := range p.bags {
		if idx == 9 {
			if len(bag.Nodes) > 1 {
				pe := New(arena, .85, ExtendedComparator)
				//pe := New(arena, .129, ExtendedComparator)
				for _, n := range bag.Nodes {
					pe.Next(n)
				}

				for _, bag := range pe.bags {
					if len(bag.Nodes) > 1 {
						for _, n := range bag.Nodes {
							bag.Volume += n.Volume
						}
						finalBags = append(finalBags, bag)

					}
				}
			}
		}
	}

	// this one is optional
	sort.Slice(finalBags, func(i, j int) bool {
		return finalBags[i].Volume > finalBags[j].Volume
	})

	bagGroups := groupBags(arena, finalBags)
	sort.Slice(bagGroups, func(i, j int) bool {
		return bagGroups[i].Volume > bagGroups[j].Volume
	})

	// transpose
	batches := make([][][]*classify.Node, len(bagGroups))
	for i, g := range bagGroups {
		batches[i] = transpose(g)
	}
	return batches, nil

	//return finalBags[rank].Nodes, nil
}

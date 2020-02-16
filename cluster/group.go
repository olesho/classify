package cluster

import (
	classify "github.com/olesho/class"
	"sort"
)

type BagGroup struct {
	Bags   []Cluster
	Volume float64
}

func groupBags(a *classify.Arena, bags []Cluster) []*BagGroup {
	for i1, bag1 := range bags {
		if i1+1 < len(bags) {
			for i2, bag2 := range bags[i1+1:] {
				if len(bag1.Members) == len(bag2.Members) {
					if belongs(a, bag1.Members, bag2.Members) {
						// remove bag2
						bags = append(bags[:i2], bags[i2+1:]...)
					} else if belongs(a, bag2.Members, bag1.Members) {
						// remove bag1
						bags = append(bags[:i1], bags[i1+1:]...)
					}
				}
			}
		}
	}

	for _, bag := range bags {
		sort.Slice(bag.Members, func(i, j int) bool {
			return bag.Members[i].Id < bag.Members[j].Id
		})
	}

	groups := []*BagGroup{}
	for _, bag1 := range bags {
		groups = checkNextIntersectionStrict(a, groups, bag1)
	}

	return groups
}

func checkNextIntersectionStrict(a *classify.Arena, groups []*BagGroup, bag1 Cluster) []*BagGroup {
	for _, g := range groups {
		if len(bag1.Members) == len(g.Bags[0].Members) {
			for i2, bag2 := range g.Bags {
				if intersects(a, bag2.Members, bag1.Members) {
					// add bag1
					g.Bags = append(g.Bags[:i2], append([]Cluster{bag1}, g.Bags[i2:]...)...)
					g.Volume += bag1.Volume()
					return groups
				} else if intersects(a, bag1.Members, bag2.Members) {
					// add bag1
					g.Bags = append(g.Bags, bag1)
					g.Volume += bag1.Volume()
					return groups
				}
			}
		}
	}
	groups = append(groups, &BagGroup{Bags: []Cluster{bag1}, Volume: bag1.Volume()})
	return groups
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

package sequence

import (
	"github.com/olesho/classify/arena"
	"sort"
)

type ClusterGroup struct {
	Clusters        []Table
	Volume          float32
	WholesomeVolume float32
	GroupVolume     float32
	Size            int
}

func groupClusters(a *arena.Arena, bags []Table) []*ClusterGroup {
	// order IDs to find intersections further
	for idx := range bags {
		sort.Slice(bags[idx].Members, func(i, j int) bool {
			return bags[idx].Members[i].Id < bags[idx].Members[j].Id
		})
	}

	var filteredBags []Table
	var fill = func(bag Table) {
		alreadyReplaced := false
		for filteredIdx := 0; filteredIdx < len(filteredBags); filteredIdx++ {
			filtered := filteredBags[filteredIdx]
			if len(bag.Members) == len(filtered.Members) {
				if belongs(a, bag.Members, filtered.Members) {
					return
				} else if belongs(a, filtered.Members, bag.Members) {
					if alreadyReplaced {
						// delete
						filteredBags = append(filteredBags[:filteredIdx], filteredBags[filteredIdx+1:]...)
						filteredIdx--
					} else {
						// replace
						filteredBags[filteredIdx] = bag
						alreadyReplaced = true
					}
				}
			}
		}
		if !alreadyReplaced {
			filteredBags = append(filteredBags, bag)
		}
	}

	for _, bag := range bags {
		fill(bag)
	}

	var groups []*ClusterGroup
	for _, bag := range filteredBags {
		groups = checkNextIntersectionStrict(a, groups, bag)
	}

	return groups
}

func checkNextIntersectionStrict(a *arena.Arena, groups []*ClusterGroup, cluster1 Table) []*ClusterGroup {
	for _, g := range groups {
		if len(cluster1.Members) == len(g.Clusters[0].Members) {
			for i2, cluster2 := range g.Clusters {
				if intersects(a, cluster2.Members, cluster1.Members) {
					// add cluster1
					g.Clusters = append(g.Clusters[:i2], append([]Table{cluster1}, g.Clusters[i2:]...)...)
					g.Volume += cluster1.Volume
					//					g.WholesomeVolume += cluster1.WholesomeVolume
					g.GroupVolume += cluster1.TemplateVolume()
					return groups
				} else if intersects(a, cluster1.Members, cluster2.Members) {
					// add cluster1
					g.Clusters = append(g.Clusters, cluster1)
					g.Volume += cluster1.Volume
					//					g.WholesomeVolume += cluster1.WholesomeVolume
					g.GroupVolume += cluster1.TemplateVolume()
					return groups
				}
			}
		}
	}
	groups = append(groups, &ClusterGroup{
		Clusters: []Table{cluster1},
		Volume:   cluster1.Volume,
		//		WholesomeVolume: cluster1.WholesomeVolume,
		GroupVolume: cluster1.TemplateVolume(),
		Size:        len(cluster1.Members),
	})
	return groups
}

func belongs(a *arena.Arena, nodes, toNodes []*arena.Node) bool {
	for i := range nodes {
		if !a.HasParent(nodes[i].Id, toNodes[i].Id) {
			return false
		}
	}
	return true
}

func intersects(a *arena.Arena, ns1, before []*arena.Node) bool {
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

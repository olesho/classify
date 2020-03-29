package cluster

import (
	"fmt"
	"github.com/olesho/classify"
	"sort"
)

const MAX_TRIES = 4

type Cluster struct {
	Members []*classify.Node
	Rate float64
}

func (c *Cluster) Volume() float64 {
	return float64(len(c.Members)) * c.Rate
}

type idxCluster struct {
	 matrix *RateMatrix
	 members  []int
	 rate float64
}

func (c *idxCluster) Volume() float64 {
	return float64(len(c.members)) * c.rate
}

func (c *idxCluster) toCluster(arena *classify.Arena) Cluster {
	result := Cluster{
		Members: make([]*classify.Node, len(c.members)),
		Rate: c.rate,
	}
	for i, memberIdx :=  range c.members {
		result.Members[i] = arena.List[memberIdx]
	}
	return result
}

func (c *idxCluster) rateCandidate(candidateIdx int) float64 {
	lowestVal := c.matrix.Cmp(c.members[0], candidateIdx)
	for _, memberIdx := range c.members[1:] {
		v := c.matrix.Cmp(memberIdx, candidateIdx)
		if v < lowestVal {
			lowestVal = v
		}
	}
	return lowestVal
}

func (c *idxCluster) hasIndex(idx int) bool {
	for _, i := range c.members {
		if i == idx {
			return true
		}
	}
	return false
}

func isin(val int, arr []int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func (c *idxCluster) nextCandidate(excluded ... int) (float64, int) {
	maxCandidateRate := .0
	maxCandidateIdx := -1
	for _, memberIdx := range c.members {
		for candidateIndex, val := range c.matrix.Rows[memberIdx] {

			// since only half table filled
			if candidateIndex < memberIdx {
				val = c.matrix.Cmp(candidateIndex, memberIdx)
			}

			if val > 0 && !c.matrix.RowExcluded[candidateIndex] && !isin(candidateIndex, excluded) {
				if !c.hasIndex(candidateIndex) {
					rate := c.rateCandidate(candidateIndex)
					if rate > maxCandidateRate {
						maxCandidateRate = rate
						maxCandidateIdx = candidateIndex
					}
				}
			}
		}
	}
	return maxCandidateRate, maxCandidateIdx
}

func (c *idxCluster) next() (*idxCluster, bool) {
	clone := &idxCluster{
		matrix: c.matrix,
		members: make([]int, len(c.members)),
		rate: c.rate,
	}
	copy(clone.members, c.members)

	excluded := make([]int, 0)
	for i := 0; i < MAX_TRIES; i++ {
		rate, idx := clone.nextCandidate(excluded...)
		if idx > -1 {
			clone.rate = rate
			clone.members = append(clone.members, idx)
			excluded = append(excluded, idx)
			if clone.Volume() > c.Volume() {
				for _, excludeIdx := range excluded {
					c.matrix.ExcludeCols(excludeIdx)
					c.matrix.ExcludeRows(excludeIdx)
				}
				return clone, true
			}
		} else {
			break
		}
	}
	return nil, false
}

func (c *idxCluster) tryAdd(candidateRate float64, candidateIndex int) bool {
	if c.Volume() < candidateRate * float64(len(c.members) + 1) {
		c.rate = candidateRate
		c.members = append(c.members, candidateIndex)
		return true
	}
	return false
}

func Extract2(arena *classify.Arena) *Matrix {
	s := NewDefaultComparator(arena)
	matrix := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float64 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	clusters := make([]Cluster, 0)
	for {
		maxRate, maxi, maxj := matrix.Max()
		if maxi < 0 {
			break
		}
		matrix.ExcludeCols(maxi)
		matrix.ExcludeRows(maxi)
		matrix.ExcludeCols(maxj)
		matrix.ExcludeRows(maxj)

		icluster := idxCluster{
			matrix: matrix,
			members: []int{maxi, maxj},
			rate: maxRate,
		}
		
		for newCluster, ok := icluster.next(); ok; newCluster, ok = icluster.next() {
			icluster = *newCluster
		}

		cluster := icluster.toCluster(arena)
		clusters = append(clusters, cluster)
	}

	// this one is optional
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Volume() > clusters[j].Volume()
	})

	bagGroups := groupBags(s.arena, clusters)
	sort.Slice(bagGroups, func(i, j int) bool {
		return bagGroups[i].Volume > bagGroups[j].Volume
	})

	// transpose
	rm := &Matrix{Arena: s.arena}
	rm.Matrix = make([][]Row, len(bagGroups))
	for i, g := range bagGroups {
		rm.Matrix[i] = transpose(g)
	}
	return rm

}

func Extract3(arena *classify.Arena) *Matrix {
	s := NewDefaultComparator(arena)
	matrix := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float64 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	//icluster := idxCluster{
	//	matrix: matrix,
	//	members: []int{1011, 1484, 1028, 1149, 1467, 1365, 1416, 1192, 1450, 1382, 1071, 1132, 1088, 1433, 1501, 1209, 1105, 1166, 1226, 1252, 1045, 1399},
	//	rate: 4.5914840714840714,
	//}

	icluster := idxCluster{
		matrix: matrix,
		members: []int{1011, 1484, 1028, 1149, 1467, 1365, 1416, 1192, 1450, 1382, 1071, 1132, 1088, 1433, 1501, 1209, 1105, 1166, 1226, 1252, 1045, 1399, 890, 907, 924},
		rate: 4.175515275959445,
	}

	for _, m := range icluster.members {
		matrix.ExcludeCols(m)
		matrix.ExcludeRows(m)
	}

	newCluster, val := icluster.next()
	//val := icluster.rateCandidate(924)
	//val := icluster.rateCandidate(890)
	fmt.Println("right rate:", val)
	fmt.Println(icluster.Volume(), "vs", newCluster.Volume())

	return nil
}

func Extract(arena *classify.Arena) *Matrix {
	s := NewDefaultComparator(arena)
	matrix := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float64 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	clusters := []Cluster{}
	for {
		_, maxi, maxj := matrix.Max()
		if maxi < 0 {
			break
		}

		cluster := Cluster{
			Members: []*classify.Node{arena.List[maxi]},
		}
		matrix.ExcludeRows(maxi)
		matrix.ExcludeCols(maxi)

		candidates := make([]Cell, 0)

		for j, r := range matrix.Rows[maxi] {
			if r > 0 && !matrix.RowExcluded[j] {
				candidates = append(candidates, Cell{j, r})
			}
		}

		for i, row := range matrix.Rows {
			if row[maxj] > 0 && !matrix.ColExcluded[maxj] {
				if i < maxi {
					candidates = append(candidates, Cell{i, row[maxj]})
				}
			}
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Rate > candidates[j].Rate
		})

		for _, candidate := range candidates {
			if len(cluster.Members) == 1 {
				cluster.Rate = candidate.Rate
				cluster.Members = append(cluster.Members, arena.List[candidate.Index])
				matrix.ExcludeCols(candidate.Index)
				matrix.ExcludeRows(candidate.Index)
			} else {
				if cluster.Volume() < candidate.Rate * float64(len(cluster.Members) + 1) {
					cluster.Rate = candidate.Rate
					cluster.Members = append(cluster.Members, arena.List[candidate.Index])
					matrix.ExcludeCols(candidate.Index)
					matrix.ExcludeRows(candidate.Index)
				} else {
					break
				}
			}
		}

		clusters = append(clusters, cluster)
	}

	// this one is optional
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Volume() > clusters[j].Volume()
	})

	//for i, c := range clusters {
	//	fmt.Println(i, arena.Chain(c.Members[0].Id, 0).XPath())
	//}

	bagGroups := groupBags(s.arena, clusters)
	sort.Slice(bagGroups, func(i, j int) bool {
		return bagGroups[i].Volume > bagGroups[j].Volume
	})

	for _, g := range bagGroups {
		fmt.Printf("len: %v, volume: %v, path: %v\n", len(g.Bags[0].Members), g.Volume, s.arena.XPath(g.Bags[0].Members[0].Id, 0))
	}

	// transpose
	rm := &Matrix{Arena: s.arena}
	rm.Matrix = make([][]Row, len(bagGroups))
	for i, g := range bagGroups {
		rm.Matrix[i] = transpose(g)
	}
	return rm

}

type Cell struct {
	Index int
	Rate float64
}
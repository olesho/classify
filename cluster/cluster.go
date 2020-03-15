package cluster

import (
	"fmt"
	"github.com/olesho/classify"
	"sort"
)

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

func (c *idxCluster) rateCandidate(candidateIdx int) float64 {
	lowestVal := c.matrix.Rows[c.members[0]][candidateIdx]
	for _, memberIdx := range c.members[1:] {
		if c.matrix.Rows[memberIdx][candidateIdx] < lowestVal {
			lowestVal = c.matrix.Rows[memberIdx][candidateIdx]
		}
	}
	return lowestVal
}

func (c *idxCluster) nextCandidate() bool {
	var maxCandidateRate float64
	maxCandidateIdx  := -1
 	for _, memberIdx := range c.members {
		for candidateIndex, val := range c.matrix.Rows[memberIdx] {
			if val > 0 {
				rate := c.rateCandidate(candidateIndex)
				if rate > maxCandidateRate {
					maxCandidateRate = rate
					maxCandidateIdx = candidateIndex
				}
			}
		}
	}

 	if maxCandidateIdx > -1 {
		if c.Volume() < maxCandidateRate * float64(len(c.members) + 1) {
			c.rate = maxCandidateRate
			c.members = append(c.members, maxCandidateIdx)
		}
	}
	return false
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

		candidates := []Cell{}

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
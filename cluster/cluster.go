package cluster

import (
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

//type Engine struct {
//	comparator Comparator
//	arena *classify.Arena
//}
//
//func New(c Comparator) *Engine {
//	return &Engine{}
//}
//
//func (e *Engine) Cmp(c1, c2 *Cluster) float64 {
//	sum := .0
//	for _, m1 := range c1.Members {
//		for _, m2 := range c2.Members {
//			sum += e.comparator.Cmp(m1, m2)
//		}
//	}
//	return sum/(float64(len(c1.Members)*len(c2.Members)))
//}

func Clusterize(arena *classify.Arena) *Matrix {
	s := NewDefaultComparator(arena)
	m := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float64 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	clusters := []Cluster{}
	for {
		_, maxi, maxj := m.Max()
		if maxi < 0 {
			break
		}

		cluster := Cluster{
			Members: []*classify.Node{arena.List[maxi]},
		}
		m.OffRows(maxi)
		m.OffCols(maxi)

		candidates := []Cell{}

		row := m.Rows[maxi]
		for j, r := range row {
			if r > 0 && !m.RowFlags[j] {
				candidates = append(candidates, Cell{j, r})
			}
		}

		for i, row := range m.Rows {
			if row[maxj] > 0 && !m.ColFlags[maxj] {
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
				m.OffCols(candidate.Index)
				m.OffRows(candidate.Index)
			} else {
				if cluster.Volume() < candidate.Rate * float64(len(cluster.Members) + 1) {
					cluster.Rate = candidate.Rate
					cluster.Members = append(cluster.Members, arena.List[candidate.Index])
					m.OffCols(candidate.Index)
					m.OffRows(candidate.Index)
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

type Cell struct {
	Index int
	Rate float64
}
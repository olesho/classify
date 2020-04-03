package cluster

import (
	"github.com/olesho/classify"
	"sort"
)

const MAX_TRIES = 4

type Cluster struct {
	Members []*classify.Node
	Rate float64
	Volume float64
	WholesomeVolume float64
}

type idxCluster struct {
	 arena *classify.Arena
	 matrix *RateMatrix
	 members  []int
	 rate float64
}

func (c *idxCluster) Volume() float64 {
	return float64(len(c.members)) * c.rate
}

func (c *idxCluster) toCluster(arena *classify.Arena, matrix *RateMatrix) Cluster {
	result := Cluster{
		Members: make([]*classify.Node, len(c.members)),
		Rate: c.rate,
	}

	// total volume derived from the least volume (smallest intersection)
	smallestVolume := arena.Get(c.members[0]).Volume
	for i, memberIdx :=  range c.members {
		result.Members[i] = arena.Get(memberIdx)
		if arena.Get(memberIdx).Volume < smallestVolume {
			smallestVolume = arena.Get(memberIdx).Volume
		}
	}
	result.Volume = smallestVolume * result.Rate * float64(len(c.members))

	result.WholesomeVolume = wholesomeVolume(arena, matrix, c.members)

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
		arena: c.arena,
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

func Extract(arena *classify.Arena) *Matrix {
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
			arena: arena,
			matrix: matrix,
			members: []int{maxi, maxj},
			rate: maxRate,
		}

		for newCluster, ok := icluster.next(); ok; newCluster, ok = icluster.next() {
			icluster = *newCluster
		}

		cluster := icluster.toCluster(arena, matrix)
		clusters = append(clusters, cluster)
	}

	// this one is optional
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Rate > clusters[j].Rate
	})

	clusterGroups := groupClusters(s.arena, clusters)
	sort.Slice(clusterGroups, func(i, j int) bool {
		//return clusterGroups[i].Volume > clusterGroups[j].Volume
		if clusterGroups[i].WholesomeVolume == clusterGroups[j].WholesomeVolume {
			return clusterGroups[i].Size > clusterGroups[j].Size
		}
		return  clusterGroups[i].WholesomeVolume > clusterGroups[j].WholesomeVolume
	})

	// transpose
	rm := &Matrix{
		Matrix: make([]Series, len(clusterGroups)),
	}
	for i, g := range clusterGroups {
		rm.Matrix[i] = Series{
			Matrix: transpose(g),
			Arena: arena,
			Volume: g.Volume,
			WholesomeVolume: g.WholesomeVolume,
			Size: g.Size,
		}
	}
	return rm

}

type Cell struct {
	Index int
	Rate float64
}
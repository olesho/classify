package cluster

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
)

const MAX_TRIES = 4

const (
	NotField = iota
	TextField
	LinkField
	ImageField
)

type Field struct {
	Type    int
	Content []string
}

func (f Field) String() string {
	s := ""
	switch f.Type {
	case TextField:
		s += "type:text\n"
	case LinkField:
		s += "type:link\n"
	case ImageField:
		s += "type:image\n"
	}
	for i, f := range f.Content {
		s += fmt.Sprintf("%v: %v\n", i, f)
	}
	return s
}

type Cluster struct {
	Arena         *arena.Arena
	TemplateArena *arena.Arena
	Members       []*arena.Node
	Rate          float64
	Volume        float64
	Table         []Field
}

type idxCluster struct {
	arena   *arena.Arena
	matrix  ComparableList
	members []int
	rate    float64
}

func (c *idxCluster) Volume() float64 {
	return float64(len(c.members)) * c.rate
}

func (c *idxCluster) toCluster(a *arena.Arena, matrix ComparableList) Cluster {
	templateArena := MergeAll(a, matrix, c.members)
	result := Cluster{
		Arena:         c.arena,
		TemplateArena: templateArena,
		Members:       make([]*arena.Node, len(c.members)),
		Rate:          c.rate,
	}

	result.Table = result.WholesomeGroupTable()

	// total volume derived from the least volume (smallest intersection)
	//smallestVolume := GetVolume(arena.Get(c.members[0]))
	for i, memberIdx := range c.members {
		result.Members[i] = a.Get(memberIdx)
		//if GetVolume(arena.Get(memberIdx)) < smallestVolume {
		//	smallestVolume = GetVolume(arena.Get(memberIdx))
		//}
	}
	//result.Volume = smallestVolume * result.Rate * float64(len(c.members))

	//result.WholesomeVolume = wholesomeVolume(arena, matrix, c.members)

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

func (c *idxCluster) nextCandidate(excluded ...int) (float64, int) {
	maxCandidateRate := .0
	maxCandidateIdx := -1
	for _, memberIdx := range c.members {
		candidateIndex, vals := c.matrix.Candidates(memberIdx)
		for _, val := range vals {

			// since only half table filled
			if candidateIndex < memberIdx {
				val = c.matrix.Cmp(candidateIndex, memberIdx)
			}

			if val > 0 && !c.matrix.IsExcluded(candidateIndex) && !isin(candidateIndex, excluded) {
				if !c.hasIndex(candidateIndex) {
					rate := c.rateCandidate(candidateIndex)
					if rate > maxCandidateRate {
						maxCandidateRate = rate
						maxCandidateIdx = candidateIndex
					}
				}
			}
			candidateIndex++
		}
	}
	return maxCandidateRate, maxCandidateIdx
}

func (c *idxCluster) next() (*idxCluster, bool) {
	clone := &idxCluster{
		arena:   c.arena,
		matrix:  c.matrix,
		members: make([]int, len(c.members)),
		rate:    c.rate,
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
					c.matrix.Exclude(excludeIdx)
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
	if c.Volume() < candidateRate*float64(len(c.members)+1) {
		c.rate = candidateRate
		c.members = append(c.members, candidateIndex)
		return true
	}
	return false
}

// Extract gets all clusters from given arena tree
func Extract(arena *arena.Arena) *Matrix {
	Init(arena)
	s := NewDefaultComparator(arena)

	// trace
	fmt.Printf("arena length: %v\n", len(arena.List))

	matrix := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float64 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	// trace
	fmt.Println("comparison matrix created:", time.Since(matrix.startedAt).Seconds(), "seconds")
	matrix.startedAt = time.Now()

	clusters := make([]Cluster, 0)
	for {
		maxRate, maxi, maxj := matrix.Max()
		if maxi < 0 {
			break
		}
		matrix.Exclude(maxi)
		matrix.Exclude(maxj)

		icluster := idxCluster{
			arena:   arena,
			matrix:  matrix,
			members: []int{maxi, maxj},
			rate:    maxRate,
		}

		// more complex and slow version
		//for newCluster, ok := icluster.next(); ok; newCluster, ok = icluster.next() {
		//	icluster = *newCluster
		//}

		// more simple and fast
		for nextVal, nextIndex := icluster.nextCandidate(); nextIndex > -1; nextVal, nextIndex = icluster.nextCandidate() {
			if !icluster.tryAdd(nextVal, nextIndex) {
				break
			}
			matrix.Exclude(nextIndex)
		}

		cluster := icluster.toCluster(arena, matrix)
		clusters = append(clusters, cluster)

		// trace
		// fmt.Println("cluster added", time.Since(matrix.startedAt).Seconds(), "seconds")
		// matrix.startedAt = time.Now()
		// fmt.Println("excluded:", matrix.excludedCount)
	}

	// this one is optional
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Rate > clusters[j].Rate
	})

	clusterGroups := groupClusters(s.arena, clusters)
	sort.Slice(clusterGroups, func(i, j int) bool {
		//return clusterGroups[i].Volume > clusterGroups[j].Volume

		//if clusterGroups[i].WholesomeVolume == clusterGroups[j].WholesomeVolume {
		//	return clusterGroups[i].Size > clusterGroups[j].Size
		//}
		//return  clusterGroups[i].WholesomeVolume > clusterGroups[j].WholesomeVolume

		return clusterGroups[i].GroupVolume > clusterGroups[j].GroupVolume
	})

	// transpose
	rm := &Matrix{
		Matrix: make([]Series, len(clusterGroups)),
	}
	for i, g := range clusterGroups {
		rm.Matrix[i] = Series{
			Matrix: transpose(g),
			Arena:  arena,
			Group:  g,
		}
	}
	return rm

}

// ExtractOptimized gets all clusters from given arena tree
func ExtractOptimized(arena *arena.Arena) *Matrix {
	Init(arena)
	s := NewDefaultComparator(arena)

	// trace
	startedAt := time.Now()
	fmt.Printf("arena length: %v\n", len(arena.List))

	matrix := NewOptimizedRateMatrixAsync(len(arena.List), len(arena.List), 4, func(i, j int) float64 {
		// matrix := NewOptimizedRateMatrix(len(arena.List), len(arena.List), func(i, j int) float64 {
		// matrix := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float64 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	// trace
	fmt.Println("comparison matrix created:", time.Since(startedAt).Seconds(), "seconds")
	startedAt = time.Now()

	clusters := make([]Cluster, 0)
	for {
		// trace
		// clusterStartedAt := time.Now()

		maxRate, maxi, maxj := matrix.Max()
		if maxi < 0 {
			break
		}
		matrix.Exclude(maxi)
		matrix.Exclude(maxj)

		icluster := idxCluster{
			arena:   arena,
			matrix:  matrix,
			members: []int{maxi, maxj},
			rate:    maxRate,
		}

		// more complex and slow version
		//for newCluster, ok := icluster.next(); ok; newCluster, ok = icluster.next() {
		//	icluster = *newCluster
		//}

		// more simple and fast
		for nextVal, nextIndex := icluster.nextCandidate(); nextIndex > -1; nextVal, nextIndex = icluster.nextCandidate() {
			if !icluster.tryAdd(nextVal, nextIndex) {
				break
			}
			matrix.Exclude(nextIndex)
		}

		cluster := icluster.toCluster(arena, matrix)
		clusters = append(clusters, cluster)

		// trace
		// fmt.Println("cluster added:", time.Since(clusterStartedAt).Seconds(), "seconds")
	}

	// trace
	fmt.Println("all clusters done:", time.Since(startedAt).Seconds(), "seconds")
	startedAt = time.Now()

	// this one is optional
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Rate > clusters[j].Rate
	})

	clusterGroups := groupClusters(s.arena, clusters)
	sort.Slice(clusterGroups, func(i, j int) bool {
		//return clusterGroups[i].Volume > clusterGroups[j].Volume

		//if clusterGroups[i].WholesomeVolume == clusterGroups[j].WholesomeVolume {
		//	return clusterGroups[i].Size > clusterGroups[j].Size
		//}
		//return  clusterGroups[i].WholesomeVolume > clusterGroups[j].WholesomeVolume

		return clusterGroups[i].GroupVolume > clusterGroups[j].GroupVolume
	})

	// transpose
	rm := &Matrix{
		Matrix: make([]Series, len(clusterGroups)),
	}
	for i, g := range clusterGroups {
		rm.Matrix[i] = Series{
			Matrix: transpose(g),
			Arena:  arena,
			Group:  g,
		}
	}

	// trace
	fmt.Println("all clusters sorted:", time.Since(startedAt).Seconds(), "seconds")

	return rm
}

type Cell struct {
	Index int
	Rate  float64
}

func (c *Cluster) TemplateVolume() float64 {
	vol := .0
	for _, row := range c.Table {
		switch row.Type {
		case TextField:
			vol += textsVolume(row.Content)
		case LinkField:
			vol += linksVolume(row.Content)
		case ImageField:
			vol += imgsVolume(row.Content)
		}
	}
	return vol
}

func uniform(strs []string) bool {
	for _, s := range strs[1:] {
		if s != strs[0] {
			return false
		}
	}
	return true
}

func textsVolume(strs []string) float64 {
	r := .0
	for _, s := range strs {
		r += float64(len(s))
	}
	return r

	//smallest := float64(len(strs[0]))
	//for _, s := range strs[1:] {
	//	val := float64(len(s))
	//	if val < smallest {
	//		smallest = val
	//	}
	//}
	//return smallest * float64(len(strs))
}

func linksVolume(strs []string) float64 {
	r := .0
	for _, s := range strs {
		if len(s) > 0 {
			r += 0.1
		}
	}
	return r
}

func imgsVolume(strs []string) float64 {
	r := .0
	for _, s := range strs {
		if len(s) > 0 {
			r += 0.1
		}
	}
	return r
}

func (c *Cluster) WholesomeGroupTable() []Field {
	result := make([]Field, 0)
	for _, n := range c.TemplateArena.List {
		if _, fieldType := WholesomeInfo(n); fieldType != NotField {
			ids := n.Ext.(*Additional).GroupIds
			if len(ids) == len(c.Members) {
				if values := extractFields(c.Arena, ids, fieldType); values != nil {
					result = append(result, *values)
				}
			}
		}
	}
	return result
}

func extractFields(arena *arena.Arena, ids []int, fieldType int) *Field {
	values := &Field{}
	values.Content = make([]string, len(ids))
	for i, id := range ids {
		values.Content[i], values.Type = WholesomeInfo(arena.Get(id))
		if values.Type != fieldType {
			return nil
		}
	}
	if !uniform(values.Content) {
		return values
	}
	return nil
}

func WholesomeInfo(n *arena.Node) (string, int) {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data), TextField
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				return attr.Val, ImageField
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return attr.Val, LinkField
			}
		}
	}
	return "", NotField
}

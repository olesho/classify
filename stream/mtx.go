package stream

import (
	"encoding/gob"
	"fmt"
	"github.com/olesho/classify/arena"
	"os"
	"sort"
	"sync"
)

type Mtx struct {
	Indexes []int
	Values [][]float32
	mutex sync.Mutex

	//_s *Storage
}

func init() {
	gob.Register(&Mtx{})
}

func (s *Storage) NewMtx(initialIndexes ...int) *Mtx {
	if len(initialIndexes) > 0 {
		m := &Mtx{
			Indexes: initialIndexes,
			Values: make([][]float32, len(initialIndexes)),
			mutex: sync.Mutex{},
		}
		for _, idx := range initialIndexes {
			s.NodeToCluster[idx] = m
		}
		return m
	}
	return &Mtx{
		Indexes: make([]int, 0),
		Values: make([][]float32, 0),
		mutex: sync.Mutex{},
	}
}

func (m *Mtx) Find(idx1, idx2 int) float32 {
	i, j := -1, -1
	for n, idx := range m.Indexes {
		if idx == idx1 {
			i = n
			break
		}
	}
	for n, idx := range m.Indexes {
		if idx == idx2 {
			j = n
			break
		}
	}

	//if i == -1 || j == -1 {
	//	fmt.Printf("indexes not found: %vx%v\n", idx1, idx2)
	//	fmt.Printf("indexes belong to same group: %v\n", m._s.NodeToCluster[idx1] == m._s.NodeToCluster[idx1])
	//	fmt.Printf("actual group: %v\n", m._s.NodeToCluster[idx1] == m)
	//	fmt.Printf("group contain: %v and %v\n", m._s.NodeToCluster[idx1].HasIdx(idx1), m._s.NodeToCluster[idx1].HasIdx(idx2))
	//	return 0
	//}

	return m.Get(i, j)
}

// Get returns similarity by given local indexes
func (m *Mtx) Get(i, j int) float32 {
	if i < j {
		diff := j - i - 1
		if diff < len(m.Values[i]) {
			return m.Values[i][diff]
		}
		return 0
	}
	if i == j {
		return 0
	}
	diff := i - j - 1
	if diff < len(m.Values[j]) {
		return m.Values[j][diff]
	}
	return 0
}

func (m *Mtx) FindIdx(idx int) int {
	for i, nextIdx := range m.Indexes {
		if idx == nextIdx {
			return i
		}
	}
	return -1
}

// Clone returns cluster matrix copy
func (m *Mtx) Clone() *Mtx {
	c := &Mtx{
		Indexes:      make([]int, len(m.Indexes)),
		Values:       make([][]float32, len(m.Values)),
	}
	copy(c.Indexes, m.Indexes)
	for i := range c.Values {
		c.Values[i] = make([]float32, len(m.Values[i]))
		copy(c.Values[i], m.Values[i])
	}
	return c
}

// Clone returns cluster matrix copy
func (m *Mtx) Equal(m2 *Mtx) bool {
	if len(m.Indexes) != len(m2.Indexes) {
		return false
	}
	size := len(m.Indexes)
	indexes1 := make([]int, size)
	indexes2 := make([]int, size)
	copy(indexes1, m.Indexes)
	copy(indexes2, m2.Indexes)
	sort.Ints(indexes1)
	sort.Ints(indexes2)

	for i := range indexes1 {
		if indexes1[i] != indexes2[i] {
			return false
		}
	}
	return true
}

// Clusters generates lvl2 clusters from matrix
func (m *Mtx) GenerateClusters() (clusters []*Cluster) {
	c := m.Clone()
	for {
		maxi, maxj, maxRate := c.max()
		if maxi < 0 {
			break
		}

		c.Values[maxi][maxj-maxi-1] = 0
		cluster := &Cluster{
			Indexes: []int{maxi, maxj},
			Rate:    maxRate,
		}

		for nextVal, nextIndex := c.nextCandidate(cluster.Indexes); nextIndex > -1; nextVal, nextIndex = c.nextCandidate(cluster.Indexes) {
			if !cluster.tryAdd(nextVal, nextIndex) {
				break
			}
			for _, idx := range cluster.Indexes {
				c.Exclude(idx, nextIndex)
			}
		}
		for _, idx := range cluster.Indexes {
			c.ExcludeRow(idx)
		}

		for i, idx := range cluster.Indexes {
			cluster.Indexes[i] = c.Indexes[idx]
		}
		clusters = append(clusters, cluster)
	}
	return clusters
}

// Clusters generates lvl2 clusters from matrix
func (m *Mtx) _GenerateClusters(a *arena.Arena) (clusters []*Cluster) {
	mm := m.Clone()
	for {
		c := mm.Clone()
		maxi, maxj, maxRate := c.max()
		if maxi < 0 {
			break
		}

		c.Values[maxi][maxj-maxi-1] = 0
		cluster := &Cluster{
			Indexes: []int{maxi, maxj},
			Rate:    maxRate,
		}
		c.Exclude(maxi, maxj)

		var nextVal float32
		var nextIndex int
		var negList = make([]*Cluster, 0)
		lastOkCluster := cluster.clone()
		for nextVal, nextIndex = c.nextCandidate(cluster.Indexes); nextIndex > -1; nextVal, nextIndex = c.nextCandidate(cluster.Indexes) {
			newCluster := cluster.add(nextVal, nextIndex)
			if newCluster.Volume() >= lastOkCluster.Volume() {
				lastOkCluster = newCluster.clone()
				negList = nil
			} else {
				negList = append(negList, newCluster)
			}
			cluster = newCluster

			if len(negList) > 5 {
				cluster = lastOkCluster
				break
			}

			for _, idx := range cluster.Indexes {
				c.Exclude(idx, nextIndex)
			}
		}
		for _, idx := range cluster.Indexes {
			mm.ExcludeRow(idx)
		}

		for i, idx := range cluster.Indexes {
			cluster.Indexes[i] = c.Indexes[idx]
		}

		classes := a.Get(cluster.Indexes[0]).Classes()
		if len(classes) > 0 {
			if classes[0] == "catalog-grid__cell" {
				fmt.Println(cluster.Rate)
				//fmt.Println(len(cluster.Indexes))
				//fmt.Println(lastOkCluster.Indexes[len(lastOkCluster.Indexes)-2], "vs", lastOkCluster.Indexes[len(lastOkCluster.Indexes)-1])
				//fmt.Println(lastOkCluster.Indexes[len(lastOkCluster.Indexes)-1], "vs", m.FindIndexes(negList[0].Indexes)[len(negList[0].Indexes)-1])
			}
		}

		clusters = append(clusters, cluster)
	}
	return clusters
}


func (m *Mtx) max() (maxI, maxJ int, val float32) {
	maxI = -1
	maxJ = -1
	for i, row := range m.Values {
		for j, curVal := range row {
			if curVal > val {
				val = curVal
				maxI = i
				maxJ = i + j + 1
			}
		}
	}
	return
}

func (m *Mtx) nextCandidate(currentIndexes []int) (float32, int) {
	var maxCandidateRate float32 = .0
	maxCandidateIdx := -1

	candidateIndexes := m.candidates(currentIndexes)
	for _, candidateIndex := range candidateIndexes {
		rate := m.rateCandidate(currentIndexes, candidateIndex)
		if rate > maxCandidateRate {
			maxCandidateRate = rate
			maxCandidateIdx = candidateIndex
		}
	}
	return maxCandidateRate, maxCandidateIdx
}

func (m *Mtx) rateCandidate(indexes []int, candidateIdx int) float32 {
	var lowestVal float32
	for _, memberIdx := range indexes {
		v := m.Get(memberIdx, candidateIdx)
		if v > 0 {
			if lowestVal == 0 {
				lowestVal = v
				continue
			}
			if v < lowestVal {
				lowestVal = v
			}
		}
	}
	return lowestVal
}

func hasIndex(idx int, indexes []int) bool{
	for _, nextIdx := range indexes {
		if nextIdx == idx {
			return true
		}
	}
	return false
}

func (m *Mtx) candidates(indexes []int) (pairIdxs []int) {
	for i := range m.Indexes {
		if !hasIndex(i, indexes) {
			pairIdxs = append(pairIdxs, i)
		}
	}
	return
}


// ExcludeRow erases all values in row and column
func (m *Mtx) ExcludeRow(index int) {
	for i := range m.Values[index] {
		m.Values[index][i] = 0
	}
	for i := range m.Values {
		if i < index {
			if index-i-1 < len(m.Values[i]) {
				m.Values[i][index-i-1] = 0
			}
		}
	}
}

// Exclude single intersection
func (m *Mtx) Exclude(y, x int) {
	if y > x {
		if y-x-1 < len(m.Values[x]) {
			m.Values[x][y-x-1] = 0
			return
		}
	}
	if x < y {
		if x-y-1 > len(m.Values[y]) {
			m.Values[y][x-y-1] = 0
		}
	}
}

func (m *Mtx) save(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(m)
}

func (m *Mtx) load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(m)
}

func (c *Cluster) hasIndex(idx int) bool {
	for _, i := range c.Indexes {
		if i == idx {
			return true
		}
	}
	return false
}
package sequence

import (
	"github.com/olesho/classify/comparator"
	"runtime"
	"sync"
	"sync/atomic"
)

type StemCluster struct {
	indexes           []int
	values            [][]float32
	stemIndexes       []int
	clusters          []*CrownCluster
	strictComparator  comparator.Comparator
	elementComparator comparator.Comparator
	root              *RootCluster

	m sync.Mutex
	stemLock sync.Mutex
}

func (c *StemCluster) addWithCrown(index int) {
	var values []float32
	if len(c.indexes) >= c.root.limit {
		values = make([]float32, c.root.limit)
	} else {
		values = make([]float32, len(c.indexes))
	}

	if len(values) > 2 {
		//async
		atomicIndex := new(int32)
		*atomicIndex = -1
		wg := sync.WaitGroup{}
		wg.Add(runtime.NumCPU())
		for cpuIdx := 0; cpuIdx < runtime.NumCPU(); cpuIdx++ {
			go func() {
				for {
					valueIndex := int(atomic.AddInt32(atomicIndex, 1))
					if valueIndex >= len(values) {
						break
					}
					j := len(c.indexes) - valueIndex - 1
					values[valueIndex] = c.root.FindStem(c.indexes[j], index) + c.root.Cmp(c.indexes[j], index)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	} else {
		// sync
		for valueIndex := range values {
			j := len(c.indexes) - valueIndex - 1
			values[valueIndex] = c.root.FindStem(c.indexes[j], index) + c.root.Cmp(c.indexes[j], index)
		}
	}

	c.values = append(c.values, values)
	c.indexes = append(c.indexes, index)

	// "local" means index in stems
	localIndex := len(c.indexes) - 1
	c.expandAnyCrown(localIndex)
}

func (c *StemCluster) addWithCrownSync(index int) {
	var values []float32
	if len(c.indexes) >= c.root.limit {
		values = make([]float32, c.root.limit)
	} else {
		values = make([]float32, len(c.indexes))
	}

	// sync
	for valueIndex := range values {
		j := len(c.indexes) - valueIndex - 1
		values[valueIndex] = c.root.FindStem(c.indexes[j], index) + c.root.Cmp(c.indexes[j], index)
	}

	c.values = append(c.values, values)
	c.indexes = append(c.indexes, index)

	// "local" means index in stems
	localIndex := len(c.indexes) - 1
	c.expandAnyCrown(localIndex)
}

func (c *StemCluster) expandAnyCrown(localIndex int) {
	var maxN int = -1
	var maxLowVal float32
	//var maxAvgVal Frac32
	// find max match to existing bags
	for n, cluster := range c.clusters {
		if low, _ := cluster.Rate(localIndex); low > maxLowVal {
			maxN = n
			maxLowVal = low
			//maxAvgVal = avg
		}
	}

	// not successful putting into any existing bag
	if maxN == -1  {
		c.addNewCrown(localIndex)
	} else {
		expanded := c.clusters[maxN].ExpandBest(maxLowVal, localIndex)
		if expanded {
			squeezedLocalIdx := c.clusters[maxN].SqueezeWorst()
			if squeezedLocalIdx > -1 {
				c.expandAnyCrown(squeezedLocalIdx)
			}
		} else {
			c.addNewCrown(localIndex)
		}
	}
}

func (c *StemCluster) addNewCrown(localIndex int) {
	c.clusters = append(c.clusters, &CrownCluster{
		items: []CrownItem{{
			Index: localIndex,
			ValueSum: 0,
		}},
		rate:    0,
		stem:    c,
	})
}

func (c *StemCluster) AddAndFillMatrix(index int) bool {
	c.stemLock.Lock()
	firstIdx := c.stemIndexes[0]
	c.stemLock.Unlock()
	fitting := c.strictComparator.Cmp(firstIdx, index)
	// if element with index fits stem cluster

	if fitting > 0 {
		if firstIdx < index {
			c.root.matrix[index][firstIdx] = fitting
		} else {
			c.root.matrix[firstIdx][index] = fitting
		}

		c.m.Lock()
		for _, existingIdx:= range c.stemIndexes {
			// calculate element fits to each existing element of cluster
			val := c.elementComparator.Cmp(index, existingIdx)
			if val > 0 {
				if existingIdx < index {
					c.root.matrix[index][existingIdx] = val
				} else {
					c.root.matrix[existingIdx][index] = val
				}
			}
		}
		c.stemLock.Lock()
		c.stemIndexes = append(c.stemIndexes, index)
		c.stemLock.Unlock()
		c.m.Unlock()
		return true
	}
	return false
}

func (c *StemCluster) AddAndFillMatrixSync(index int) bool {
	firstIdx := c.stemIndexes[0]
	fitting := c.strictComparator.Cmp(firstIdx, index)
	// if element with index fits stem cluster
	if fitting > 0 {
		for _, existingIdx := range c.stemIndexes {
			// calculate element fits to each existing element of cluster
			val := c.elementComparator.Cmp(index, existingIdx)
			if val > 0 {
				if existingIdx < index {
					c.root.matrix[index][existingIdx] = val
				} else {
					c.root.matrix[existingIdx][index] = val
				}
			}
		}
		// append to cluster
		c.stemLock.Lock()
		c.stemIndexes = append(c.stemIndexes, index)
		c.stemLock.Unlock()
		return true
	}
	return false
}

func (c *StemCluster) Get(i, j int) float32 {
	if i < j {
		diff := j - i - 1
		if diff >= c.root.limit {
			return 0
		}
		return c.values[j][diff]
	} else if j < i {
		diff := i - j - 1
		if diff >= c.root.limit {
			return 0
		}
		return c.values[i][diff]
	}
	return 0
}

func (c *StemCluster) FindIdx(idx int) int {
	for i, next := range c.indexes {
		if next == idx {
			return i
		}
	}
	return -1
}

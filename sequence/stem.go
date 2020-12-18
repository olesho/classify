package sequence

import (
	"github.com/olesho/classify/comparator"
	"runtime"
	"sync"
	"sync/atomic"
)

type ending struct {
	last int
	index int
}

type StemCluster struct {
	indexes []int
	values [][]float32
	stemIndexes []int
	clusters []*CrownCluster
	strictComparator comparator.Comparator
	elementComparator comparator.Comparator
	root *RootCluster

	m sync.Mutex
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
	localIndex := len(c.indexes)-1

	var maxN int = -1
	var maxVal float32
	// find max match to existing bags
	for n, cluster := range c.clusters {
		if val := cluster.Rate(localIndex); val > maxVal {
			maxN = n
			maxVal = val
		}
	}

	// not successful putting into any existing bag
	if maxN == -1 || !c.clusters[maxN].Add(maxVal, localIndex) {
		c.clusters = append(c.clusters, &CrownCluster{
			indexes: []int{localIndex},
			rate:    1,
			stem:    c,
		})
	}
}

func (c *StemCluster) AddFirst(index int)  {
	c.stemIndexes = []int{index}
	last := c.root.Arena.Get(index).Ext.(*Additional).LastDescendant
	if index == last {
		//same as c.addWithCrown(0, index)
		c.indexes = []int{index}
		c.values = append(c.values, []float32{})
		c.clusters = append(c.clusters, &CrownCluster{
			indexes: []int{0},
			rate:    1,
			stem:    c,
		})
	}
}

func (c *StemCluster) Add(index int) bool {
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
		c.m.Lock()
		defer c.m.Unlock()
		c.stemIndexes = append(c.stemIndexes, index)
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
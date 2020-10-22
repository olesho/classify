package sequence

import (
	"github.com/olesho/classify/comparator"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

type ending struct {
	last int
	index int
	i int
}

type StemCluster struct {
	indexes []int
	values [][]float32

	endings []ending
	stemIndexes []int
	stemValues [][]float32

	clusters []*CrownCluster

	strictComparator comparator.Comparator
	elementComparator comparator.Comparator
	root *RootCluster

	m sync.Mutex
}

func (c *StemCluster) addWithCrown(indexI, index int) {
	//values := make([]float32, len(c.indexes))

	var values []float32
	if len(c.indexes) >= c.root.limit {
		values = make([]float32, c.root.limit)
	} else {
		values = make([]float32, len(c.indexes))
	}

	//async
	atomicIndex := new(int32)
	*atomicIndex = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for cpuIdx := 0; cpuIdx < runtime.NumCPU(); cpuIdx++ {
		go func() {
			for i := int(atomic.AddInt32(atomicIndex, 1)); i < len(values); i = int(atomic.AddInt32(atomicIndex, 1)) {
				values[i] = c.GetStem(i, indexI) + c.root.Cmp(c.indexes[i], index)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// sync
	//for i := range values {
	//	values[i] = c.GetStem(i, indexI) + c.root.Cmp(c.indexes[i], index)
	//}

	c.values = append(c.values, values)
	c.indexes = append(c.indexes, index)

	idx := len(c.indexes)-1
	var maxI int = -1
	var maxVal float32
	// find max match to existing bags
	for i, cluster := range c.clusters {
		if val := cluster.Rate(idx); val > maxVal {
			maxI = i
			maxVal = val
		}
	}

	// not successful putting into any existing bag
	if maxI == -1 {
		c.clusters = append(c.clusters, &CrownCluster{
			indexes: []int{index},
			rate:    1,
			stem:    c,
		})
	} else {
		c.clusters[maxI].Add(maxVal, index)
	}
}

func (c *StemCluster) Notify(index int) {
	if len(c.endings) > 0 {
		if index > c.endings[0].last {
			c.addWithCrown(c.endings[0].i, c.endings[0].index)
			c.endings = c.endings[1:]
		}
	}
}

func (c *StemCluster) AddFirst(index int) bool {
	c.stemIndexes = []int{index}
	c.stemValues = make([][]float32, 1)
	last := c.root.arena.Get(index).Ext.(*Additional).LastDescendant
	if index == last {
		//same as c.addWithCrown(0, index)
		c.clusters = append(c.clusters, &CrownCluster{
			indexes: []int{index},
			rate:    1,
			stem:    c,
		})
	} else {
		c.endings = []ending{
			{
				i:     0,
				index: index,
				last:  c.root.arena.Get(index).Ext.(*Additional).LastDescendant,
			},
		}
	}
	return true
}

func (c *StemCluster) Add(index int) bool {
	if c.strictComparator.Cmp(c.stemIndexes[0], index) > 0 {
		values := make([]float32, len(c.stemIndexes))
		for i, existingIdx := range c.stemIndexes {
			if val := c.elementComparator.Cmp(index, existingIdx); val > 0 {
				values[i] = val
			} else {
				return false
			}
		}
		c.stemIndexes = append(c.stemIndexes, index)
		c.stemValues = append(c.stemValues, values)

		c.m.Lock()
		c.endings = append(c.endings, ending{
			i:     len(c.stemIndexes)-1,
			index: index,
			last:  c.root.arena.Get(index).Ext.(*Additional).LastDescendant,
		})
		sort.Slice(c.endings, func(i, j int) bool {
			return c.endings[i].last > c.endings[j].last
		})
		c.m.Unlock()
		return true
	}
	return false
}

func (c *StemCluster) GetStem(i, j int) float32 {
	if i < j {
		return c.stemValues[j][i]
	} else if j < i {
		return c.stemValues[i][j]
	}
	return 0
}

func (c *StemCluster) FindStem(idx1, idx2 int) float32 {
	i, j := -1, -1
	for n, idx := range c.stemIndexes {
		if idx == idx1 {
			i = n
			break
		}
	}
	for n, idx := range c.stemIndexes {
		if idx == idx2 {
			j = n
			break
		}
	}
	return c.GetStem(i, j)
}

func (c *StemCluster) Get(i, j int) float32 {
	if i < j {
		if i >= len(c.values[j]) {
			return 0
		}
		return c.values[j][i]
	} else if j < i {
		if j >= len(c.values[i]) {
			return 0
		}
		return c.values[i][j]
	}
	return 0
}

func (c *StemCluster) Find(idx1, idx2 int) float32 {
	i, j := -1, -1
	for n, idx := range c.indexes {
		if idx == idx1 {
			i = n
			break
		}
	}
	for n, idx := range c.indexes {
		if idx == idx2 {
			j = n
			break
		}
	}
	return c.Get(i, j)
}
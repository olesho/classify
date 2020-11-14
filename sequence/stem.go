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
}

type StemCluster struct {
	indexes []int
	values [][]float32

	endings []ending
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

	//async
	atomicIndex := new(int32)
	*atomicIndex = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for cpuIdx := 0; cpuIdx < runtime.NumCPU(); cpuIdx++ {
		go func() {
			for valueIndex := int(atomic.AddInt32(atomicIndex, 1)); valueIndex < len(values); valueIndex = int(atomic.AddInt32(atomicIndex, 1)) {
				j := len(c.indexes) - valueIndex - 1
				values[valueIndex] = c.root.FindStem(c.indexes[j], index) + c.root.Cmp(c.indexes[j], index)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	////sync
	//for valueIndex := range values {
	//	j := len(c.indexes) - valueIndex - 1
	//	stemVal := c.root.FindStem(c.indexes[j], index)
	//	crownVal := c.root.Cmp(c.indexes[j], index)
	//	values[valueIndex] = crownVal + stemVal
	//}

	c.values = append(c.values, values)
	c.indexes = append(c.indexes, index)
	c.add(len(c.indexes)-1)
}

func (c *StemCluster) maxAttachCluster(localIndex int, exceptClusters ... int) (clusterIndex int, maxRate float32) {
	clusterIndex = -1
	// find max match to existing bags
	for n, cluster := range c.clusters {
		if isin(n, exceptClusters) {
			continue
		}

		if val := cluster.attachRate(localIndex); val > maxRate {
			clusterIndex = n
			maxRate = val
		}
	}
	return
}

func (c *StemCluster) add(localIndex int, exceptClusters ... int) {
	maxN, maxVal := c.maxAttachCluster(localIndex, exceptClusters...)
	var addedSuccessfully bool
	if maxN > -1 {
		bestMatchingCluster := c.clusters[maxN]
		attachDelta := bestMatchingCluster.attachDelta(maxVal)
		addedSuccessfully = attachDelta > 0
		if addedSuccessfully {
			bestMatchingCluster.attach(maxVal, localIndex)
			detachedCrownIndex, detachDelta, rateAfterDetachment := bestMatchingCluster.detachRate()
			if detachedCrownIndex > -1 &&
				bestMatchingCluster.indexes[detachedCrownIndex] != localIndex { // cannot detach what have just been attached
				if detachDelta > 0 {
					c.add(bestMatchingCluster.indexes[detachedCrownIndex], maxN)
					bestMatchingCluster.detach(rateAfterDetachment, detachedCrownIndex)
				} else if len(c.clusters) > 1 { // are there alternatives
					nextMaxN, maxAttachRate := c.maxAttachCluster(bestMatchingCluster.indexes[detachedCrownIndex], maxN)
					if nextMaxN > -1 {
						nextBestMatchingCluster := c.clusters[nextMaxN]
						maxAttachDelta := nextBestMatchingCluster.attachDelta(maxAttachRate)
						if (maxAttachDelta + detachDelta) > 0 {
							nextBestMatchingCluster.attach(maxAttachRate, bestMatchingCluster.indexes[detachedCrownIndex])
							bestMatchingCluster.detach(rateAfterDetachment, detachedCrownIndex)
						}
					}
				}
			}
		}
	}

	// not successful putting into any existing bag
	if maxN == -1 || !addedSuccessfully {
		c.clusters = append(c.clusters, &CrownCluster{
			indexes: []int{localIndex},
			rate:    1,
			stem:    c,
		})
	}
}

func (c *StemCluster) AddFirst(index int) bool {
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
	} else {
		c.endings = []ending{
			{
				index: index,
				last:  c.root.Arena.Get(index).Ext.(*Additional).LastDescendant,
			},
		}
	}
	return true
}

func (c *StemCluster) Add(index int) bool {
	if c.strictComparator.Cmp(c.stemIndexes[0], index) > 0 {
		for _, existingIdx := range c.stemIndexes {
			if val := c.elementComparator.Cmp(index, existingIdx); val > 0 {
				c.root.matrix[index][existingIdx] = val
			}
			// this should never happen
			//} else {
			//	return false
			//}
		}
		c.stemIndexes = append(c.stemIndexes, index)
		//c.stemValues = append(c.stemValues, values)

		c.m.Lock()
		c.endings = append(c.endings, ending{
			index: index,
			last:  c.root.Arena.Get(index).Ext.(*Additional).LastDescendant,
		})
		sort.Slice(c.endings, func(i, j int) bool {
			return c.endings[i].last > c.endings[j].last
		})
		c.m.Unlock()
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

func (c *StemCluster) Find(idx int) int {
	for i, index := range c.indexes {
		if index == idx {
			return i
		}
	}
	return -1
}
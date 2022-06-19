package sequence

import (
	"bufio"
	"fmt"
	"github.com/olesho/classify/arena"
	"github.com/olesho/classify/comparator"
	"golang.org/x/net/html"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type crownConsumer struct {
	wg                   sync.WaitGroup
	stemIndexDone        chan struct{}
	awaitingLock         sync.Mutex
	awaitingCompareCrown []compareCrown
	lastNotified         int
}

type DebugConfig struct {
	DebugMatrix bool
	DebugExpansion bool
	DebugClusterization bool
	DebugGroups bool
	TagName string
	AttrKey string
	AttrVal string
}

type RootCluster struct {
	limit int
	debug *DebugConfig

	clusters        []*StemCluster
	nodeIDToCluster []*StemCluster
	matrix          [][]float32

	Arena             *arena.Arena
	strictComparator  comparator.Comparator
	elementComparator comparator.Comparator

	consumer *crownConsumer

	m sync.Mutex
}

type compareCrown struct {
	stemCluster    *StemCluster
	index          int
	lastDescendant int
}

func NewRootCluster() *RootCluster {
	a := arena.NewArena()
	return &RootCluster{
		limit:             99999,
		Arena:             a,
		strictComparator:  comparator.NewStrictComparator(a),
		elementComparator: comparator.NewElementComparator(a),

		consumer: &crownConsumer{
			wg:            sync.WaitGroup{},
			stemIndexDone: make(chan struct{}),
			awaitingLock:  sync.Mutex{},
		},

		m: sync.Mutex{},
	}
}

func (rs *RootCluster) SetLimit(limit int) *RootCluster {
	rs.limit = limit
	return rs
}

func (rs *RootCluster) matchDebug(n *arena.Node) bool {
	if n.Type == html.ElementNode &&
		n.Data == rs.debug.TagName {
		if rs.debug.AttrKey != "" {
			if rs.debug.AttrVal != "" && n.GetAttr(rs.debug.AttrKey) == rs.debug.AttrVal {
				return true
			}
			if n.GetAttr(rs.debug.AttrKey) != "" {
				return true
			}
			return false
		}
		return true
	}
	return false
}

func (rs *RootCluster) SetDebug(config *DebugConfig) *RootCluster {
	rs.debug = config
	return rs
}

func (rs *RootCluster) newStemCluster(index int) *StemCluster {
	sc := &StemCluster{
		strictComparator:  rs.strictComparator,
		elementComparator: rs.elementComparator,
		root:              rs,
		m:                 sync.Mutex{},
		stemLock:          sync.Mutex{},
		stemIndexes:       []int{index},
	}
	//sc.AddFirst(index)
	return sc
}

// LoadFile appends HTML file content
func (rs *RootCluster) Load(r io.Reader) error {
	n, err := html.Parse(r)
	if err != nil {
		return err
	}
	rs.Arena.Append(*n)
	return nil
}

// LoadFile appends HTML file content
func (rs *RootCluster) LoadFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	if err != nil {
		return err
	}
	rs.Arena.Append(*n)
	return nil
}

// LoadString appends HTML string
func (rs *RootCluster) LoadString(str string) error {
	n, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return err
	}
	rs.Arena.Append(*n)
	return nil
}

func (rs *RootCluster) Batch() *RootCluster {
	Init(rs.Arena)
	rs.nodeIDToCluster = make([]*StemCluster, len(rs.Arena.List))
	rs.matrix = make([][]float32, len(rs.Arena.List))
	for i := range rs.matrix {
		rs.matrix[i] = make([]float32, i)
	}

	atomicIndex := new(int32)
	*atomicIndex = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for cpuIdx := 0; cpuIdx < runtime.NumCPU(); cpuIdx++ {
		go func() {
			for {
				valueIndex := int(atomic.AddInt32(atomicIndex, 1))
				if valueIndex >= len(rs.Arena.List) {
					break
				}
				rs.Process(valueIndex)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	// merge some stem cluster as previous async operation might have produced clusters of same kind
	for i, firstCluster := range rs.clusters[:len(rs.clusters)-1] {
		if firstCluster != nil {
			for j := i + 1; j < len(rs.clusters); j++ {
				secondCluster := rs.clusters[j]
				if secondCluster == nil {
					continue
				}
				secondCluster.stemLock.Lock()
				if firstCluster.AddAndFillMatrixSync(secondCluster.stemIndexes[0]) {
					rs.nodeIDToCluster[secondCluster.stemIndexes[0]] = rs.clusters[i]
					for _, nextIdx := range secondCluster.stemIndexes[1:] {
						firstCluster.AddAndFillMatrixSync(nextIdx)
						rs.nodeIDToCluster[nextIdx] = firstCluster
					}
					rs.clusters[j] = nil
				}
				secondCluster.stemLock.Unlock()
			}
		}
	}
	// remove nils
	for i := 0; i < len(rs.clusters); i++ {
		if rs.clusters[i] == nil {
			rs.clusters = append(rs.clusters[:i], rs.clusters[i+1:]...)
			i--
		} else {
			sort.Ints(rs.clusters[i].stemIndexes)
		}
	}

	sort.Slice(rs.clusters, func(i, j int) bool {
		return rs.clusters[i].stemIndexes[0] < rs.clusters[j].stemIndexes[0]
	})

	//rs.consumeNotificationsSync()
	rs.consumeNotifications()
	return rs
}

func (rs *RootCluster) BatchSync() *RootCluster {
	Init(rs.Arena)
	rs.nodeIDToCluster = make([]*StemCluster, len(rs.Arena.List))
	rs.matrix = make([][]float32, len(rs.Arena.List))
	for i := range rs.matrix {
		rs.matrix[i] = make([]float32, i)
	}
	for i := range rs.Arena.List {
		rs.Process(i)
	}
	rs.consumeNotificationsSync()
	return rs
}

func (rs *RootCluster) consumeNotifications() {
	for _, stemCluster := range rs.clusters {
		for _, index := range stemCluster.stemIndexes {
			stemCluster.addWithCrown(index)
		}

		for idx1, idx2, val := stemCluster.findCrownCandidatesToMerge();
			idx1 > -1 && idx2 > -1;
		idx1, idx2, val = stemCluster.findCrownCandidatesToMerge() {
			stemCluster.clusters[idx1].items = append(stemCluster.clusters[idx1].items, stemCluster.clusters[idx2].items...)
			stemCluster.clusters[idx1].rate = val
			stemCluster.clusters = append(stemCluster.clusters[:idx2], stemCluster.clusters[idx2+1:]...)
			for idx := stemCluster.clusters[idx1].SqueezeWorst();
				idx > -1;
			idx = stemCluster.clusters[idx1].SqueezeWorst() {}
		}

		if rs.debug != nil && rs.debug.TagName == rs.Arena.Get(stemCluster.stemIndexes[0]).Data {
			for _, c := range stemCluster.clusters {
				element := rs.Arena.Get(stemCluster.indexes[c.items[0].Index])
				fmt.Printf("clusterized (%v):%v\n", len(c.items), element)
			}
		}
	}
}

func (rs *RootCluster) consumeNotificationsSync() {
	for _, stemCluster := range rs.clusters {
		for _, index := range stemCluster.stemIndexes {
			stemCluster.addWithCrownSync(index)
		}

		for idx1, idx2, val := stemCluster.findCrownCandidatesToMerge();
		idx1 > -1 && idx2 > -1;
		idx1, idx2, val = stemCluster.findCrownCandidatesToMerge() {
			stemCluster.clusters[idx1].items = append(stemCluster.clusters[idx1].items, stemCluster.clusters[idx2].items...)
			stemCluster.clusters[idx1].rate = val
			stemCluster.clusters = append(stemCluster.clusters[:idx2], stemCluster.clusters[idx2+1:]...)
			for idx := stemCluster.clusters[idx1].SqueezeWorst();
				idx > -1;
				idx = stemCluster.clusters[idx1].SqueezeWorst() {}
		}

		if rs.debug != nil && rs.debug.TagName == rs.Arena.Get(stemCluster.stemIndexes[0]).Data {
			for _, c := range stemCluster.clusters {
				element := rs.Arena.Get(stemCluster.indexes[c.items[0].Index])
				fmt.Printf("clusterized (%v):%v\n", len(c.items), element)
			}
		}
	}
}

func (c *StemCluster) findCrownCandidatesToMerge() (int, int, float32) {
	var max float32
	maxIdx1, maxIdx2 := -1, -1
	for idx1 := range c.clusters {
		for idx2 := idx1 + 1; idx2 < len(c.clusters); idx2++ {
			val := c.evalCrownMerge(idx1, idx2)
			if val > max {
				max = val
				maxIdx1 = idx1
				maxIdx2 = idx2
			}
		}
	}
	return maxIdx1, maxIdx2, max
}

func (c *StemCluster) evalCrownMerge(idx1, idx2 int) float32 {
	c1 := c.clusters[idx1]
	c2 := c.clusters[idx2]
	var lowest float32
	for _, item1 := range c1.items {
		low, _ := c2.RateAgainst(item1.Index)
		if lowest == 0 { lowest = low }
		if low < lowest {
			lowest = low
		}
	}
	return lowest
}

func (rs *RootCluster) Process(index int) {
	// try add into one of existing bags
	var i int
	for i = 0; i < len(rs.clusters); i++ {
		rs.m.Lock()
		next := rs.clusters[i]
		rs.m.Unlock()
		if next.AddAndFillMatrix(index) {
			rs.nodeIDToCluster[index] = next
			return
		}
	}

	// not successful putting into any existing bag
	stemCluster := rs.newStemCluster(index)
	rs.nodeIDToCluster[index] = stemCluster

	rs.m.Lock()
	defer rs.m.Unlock()
	rs.clusters = append(rs.clusters, stemCluster)
	return
}

func (rs *RootCluster) Results() []*Series {
	var crownClusters = make([]*CrownCluster, 0)
	for _, stemCluster := range rs.clusters {
		for _, crownCluster := range stemCluster.clusters {
			crownCluster.resolveIndexes()
			crownClusters = append(crownClusters, crownCluster)
		}
	}

	sort.Slice(crownClusters, func(i, j int) bool {
		return len(crownClusters[i].items) > len(crownClusters[j].items)
	})

	tables := make([]Table, len(crownClusters))

	atomicIndex := new(int32)
	*atomicIndex = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for cpuIdx := 0; cpuIdx < runtime.NumCPU(); cpuIdx++ {
		go func() {
			for {
				valueIndex := int(atomic.AddInt32(atomicIndex, 1))
				if valueIndex >= len(crownClusters) {
					break
				}
				tables[valueIndex] = crownClusters[valueIndex].toTable()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	//for i, cluster := range crownClusters {
	//	tables[i] = cluster.toTable()
	//}

	clusterGroups := groupClusters(rs.Arena, tables)
	if rs.debug != nil {
		for _, g := range clusterGroups {
			for _, cluster := range g.Clusters {
				if rs.matchDebug(cluster.Members[0]) && rs.debug.DebugGroups {
					fmt.Printf("groups: %v\n", cluster.Members[0])
				}
			}
		}
	}

	// transpose
	rm := make([]*Series, len(clusterGroups))
	for i, g := range clusterGroups {
		//rm[i] = removeEqualFields(transpose(g))
		rm[i] = transpose(g)
		rm[i].Arena = rs.Arena
		rm[i].Group.Volume = rateSeries(rm[i])
	}

	sort.Slice(rm, func(i, j int) bool {
		return rm[i].Group.Volume > rm[j].Group.Volume
		//return clusterGroups[i].GroupVolume > clusterGroups[j].GroupVolume
	})

	return rm
}

func equalFields(f1, f2 []string) bool {
	if len(f1) == len(f2) {
		for i := range f1 {
			if f1[i] != f2[i] {
				return false
			}
		}
	}
	return true
}

func removeEqualFields(s *Series) *Series {
	for fieldIndex := 0; fieldIndex < len(s.TransposedValues); fieldIndex++ {
		for offset := fieldIndex + 1; offset < len(s.TransposedValues); offset++ {
			if equalFields(s.TransposedValues[fieldIndex], s.TransposedValues[offset]) {
				s.TransposedValues = append(s.TransposedValues[:offset], s.TransposedValues[offset+1:]...)
				s.TransposedNodes = append(s.TransposedNodes[:offset], s.TransposedNodes[offset+1:]...)
				offset--
			}
		}
	}
	return s
}

func rateSeries(s *Series) float32 {
	var sum float32
	for _, fields := range s.TransposedValues {
		for range fields {
			sum += 1
		}
	}
	return sum
}

func (rs *RootCluster) FindStem(idx1, idx2 int) float32 {
	if idx1 > idx2 {
		return rs.matrix[idx1][idx2]
	} else if idx2 > idx1 {
		return rs.matrix[idx2][idx1]
	}
	return 0
}

func (rs *RootCluster) FindCrown(idx1, idx2 int) float32 {
	c1, c2 := rs.nodeIDToCluster[idx1], rs.nodeIDToCluster[idx2]
	if c1 != c2 {
		return 0
	}

	i, j := -1, -1
	for n, idx := range c1.indexes {
		if idx == idx1 {
			i = n
			break
		}
	}
	for n, idx := range c1.indexes {
		if idx == idx2 {
			j = n
			break
		}
	}

	return c2.Get(i, j)
}

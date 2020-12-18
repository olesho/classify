package sequence

import (
	"bufio"
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
	wg sync.WaitGroup
	stemIndexDone chan struct{}
	awaitingLock sync.Mutex
	awaitingCompareCrown []compareCrown
	lastNotified int
}

type RootCluster struct {
	limit int

	clusters []*StemCluster
	nodeIDToCluster []*StemCluster
	matrix [][]float32

	Arena *arena.Arena
	strictComparator comparator.Comparator
	elementComparator comparator.Comparator

	consumer *crownConsumer

	m sync.Mutex
}

type compareCrown struct {
	stemCluster *StemCluster
	index int
	lastDescendant int
}

func NewRootCluster() *RootCluster {
	a := arena.NewArena()
	return &RootCluster{
		limit: 99999,
		Arena: a,
		strictComparator: comparator.NewStrictComparator(a),
		elementComparator: comparator.NewElementComparator(a),

		consumer: &crownConsumer{
			wg: sync.WaitGroup{},
			stemIndexDone: make(chan struct{}),
			awaitingLock: sync.Mutex{},
		},

		m: sync.Mutex{},
	}
}

func (rs *RootCluster) SetLimit(limit int) *RootCluster {
	rs.limit = limit
	return rs
}

func (rs *RootCluster) newStemCluster(index int) *StemCluster {
	sc := &StemCluster{
		strictComparator:  rs.strictComparator,
		elementComparator: rs.elementComparator,
		root: rs,
		m: sync.Mutex{},
	}
	sc.AddFirst(index)
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
				rs.Add(valueIndex)
			}
			wg.Done()
		}()
	}
	wg.Wait()
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
		rs.Add(i)
	}
	rs.consumeNotifications()
	return rs
}

func (rs *RootCluster) consumeNotifications() {
	for _, sc := range rs.clusters {
		for _, index  := range sc.stemIndexes {
			sc.addWithCrown(index)
		}
	}
}

func (rs *RootCluster) Add(index int) {
	// try add into one of existing bags
	var i int
	for i = 0; i < len(rs.clusters); i++ {
		if rs.clusters[i].Add(index) {
			rs.nodeIDToCluster[index] = rs.clusters[i]
			return
		}
	}

	// not successful putting into any existing bag
	stemCluster := rs.newStemCluster(index)
	rs.nodeIDToCluster[index] = stemCluster
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
	sort.Slice(crownClusters, func(i,j int) bool {
		return len(crownClusters[i].indexes) > len(crownClusters[j].indexes)
	})

	tables := make([]Table, len(crownClusters))
	for i, cluster := range crownClusters {
		tables[i] = cluster.toTable()
	}

	clusterGroups := groupClusters(rs.Arena, tables)

	// transpose
	rm := make([]*Series, len(clusterGroups))
	for i, g := range clusterGroups {
		rm[i] = removeEqualFields(transpose(g))
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
	for fieldIndex := 0; fieldIndex <  len(s.TransposedFields); fieldIndex++ {
		for offset := fieldIndex + 1; offset < len(s.TransposedFields); offset++ {
			if equalFields(s.TransposedFields[fieldIndex], s.TransposedFields[offset]) {
				s.TransposedFields = append(s.TransposedFields[:offset], s.TransposedFields[offset+1:]...)
				s.TransposedNodes = append(s.TransposedNodes[:offset], s.TransposedNodes[offset+1:]...)
				offset--
			}
		}
	}
	return s
}

func rateSeries(s *Series) float32 {
	var sum float32
	for _, fields := range s.TransposedFields {
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

//func (rs *RootCluster) String() string {
//	res := ""
//	for _, crownCluster := range rs.Results() {
//		res += fmt.Sprintln(crownCluster)
//	}
//	return res
//}
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
)

type RootCluster struct {
	limit int

	clusters []*StemCluster
	nodeIDToCluster []*StemCluster
	matrix [][]float32

	Arena *arena.Arena
	strictComparator comparator.Comparator
	elementComparator comparator.Comparator

	wg sync.WaitGroup
	notify chan [2]int
}

func NewRootCluster() *RootCluster {
	a := arena.NewArena()
	return &RootCluster{
		limit: 99999,
		Arena: a,
		strictComparator: comparator.NewStrictComparator(a),
		elementComparator: comparator.NewElementComparator(a),
		wg: sync.WaitGroup{},
		notify: make(chan [2]int),
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
	rs.consumeNotifications()
	for i := range rs.Arena.List {
		rs.Add(i)
	}
	rs.notifyAll()
	for len(rs.notify) > 0 {}
	close(rs.notify)
	rs.wg.Wait()
	return rs
}

func (rs *RootCluster) consumeNotifications() {
	rs.wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(){
			for pair := range rs.notify {
				c := rs.clusters[pair[0]]
				index := pair[1]

				c.m.Lock()
				if len(c.endings) > 0 {
					lastEndingIndex := len(c.endings)-1
					if index > c.endings[len(c.endings)-1].last {
						c.addWithCrown(c.endings[lastEndingIndex].index)
						c.endings = c.endings[:lastEndingIndex]
					}
				}
				c.m.Unlock()
			}
			rs.wg.Done()
		}()
	}
}

func (rs *RootCluster) notifyAll() {
	for i := range rs.clusters {
		rs.notify <- [2]int{i, len(rs.Arena.List)}
	}
}

func (rs *RootCluster) notifyIndex(index int) {
	for i := range rs.clusters {
		rs.notify <- [2]int{i, index}
	}
}

//func (rs *RootCluster) Rate(index int) float32 { return 0 }
func (rs *RootCluster) Add(index int) bool {
	defer rs.notifyIndex(index)

	// try add into one of existing bags
	var i int
	for i = 0; i < len(rs.clusters); i++ {
		if rs.clusters[i].Add(index) {
			rs.nodeIDToCluster[index] = rs.clusters[i]
			return true
		}
	}

	// not successful putting into any existing bag
	if i == len(rs.clusters) {
		stemCluster := rs.newStemCluster(index)
		rs.nodeIDToCluster[index] = stemCluster
		rs.clusters = append(rs.clusters, stemCluster)
	}

	return true
}

func (rs *RootCluster) Results() []*Series {
	var crownClusters = make([]*CrownCluster, 0)
	for _, stemCluster := range rs.clusters {
		for _, crownCluster := range stemCluster.clusters {
			crownCluster.extend()
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
	return c1.Get(i, j)
}

func (rs *RootCluster) String() string {
	res := ""
	for _, crownCluster := range rs.Results() {
		res += fmt.Sprintln(crownCluster)
	}
	return res
}
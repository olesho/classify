package sequence

import (
	"bufio"
	"fmt"
	"github.com/olesho/classify/arena"
	"github.com/olesho/classify/comparator"
	"golang.org/x/net/html"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type RootCluster struct {
	clusters []*StemCluster
	nodeIDToCluster []*StemCluster

	arena *arena.Arena
	strictComparator comparator.Comparator
	elementComparator comparator.Comparator
}

func NewRootCluster() *RootCluster {
	a := arena.NewArena()
	return &RootCluster{
		arena: a,
		strictComparator: comparator.NewStrictComparator(a),
		elementComparator: comparator.NewElementComparator(a),
	}
}

func (rs *RootCluster) newStemCluster(index int) *StemCluster {
	sc := &StemCluster{
		arena:             rs.arena,
		strictComparator:  rs.strictComparator,
		elementComparator: rs.elementComparator,
		root: rs,
		m: sync.Mutex{},
	}
	sc.AddFirst(index)
	return sc
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
	rs.arena.Append(*n)
	return nil
}

// LoadString appends HTML string
func (rs *RootCluster) LoadString(str string) error {
	n, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return err
	}
	rs.arena.Append(*n)
	return nil
}

func (rs *RootCluster) BatchAsync() {
	Init(rs.arena)
	rs.nodeIDToCluster = make([]*StemCluster, len(rs.arena.List))

	atomicIndex := new(int32)
	*atomicIndex = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for cpuIdx := 0; cpuIdx < runtime.NumCPU(); cpuIdx++ {
		go func() {
			for i := int(atomic.AddInt32(atomicIndex, 1)); i < len(rs.arena.List); i = int(atomic.AddInt32(atomicIndex, 1)) {
				rs.Add(i)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	rs.notifyAll()
}

func (rs *RootCluster) Batch() {
	Init(rs.arena)
	rs.nodeIDToCluster = make([]*StemCluster, len(rs.arena.List))
	for i := range rs.arena.List {
		rs.Add(i)
	}
	rs.notifyAll()
}

func (rs *RootCluster) notifyAll() {
	for _, cluster := range rs.clusters {
		cluster.Notify(len(rs.arena.List))
	}
}

//func (rs *RootCluster) Rate(index int) float32 { return 0 }
func (rs *RootCluster) Add(index int) bool {
	//if len(rs.clusters) < runtime.NumCPU() {
	//	for _, cluster := range rs.clusters {
	//		cluster.Notify(index)
	//	}
	//} else {
	atomicIndex := new(int32)
	*atomicIndex = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for cpuIdx := 0; cpuIdx < runtime.NumCPU(); cpuIdx++ {
		go func() {
			for i := int(atomic.AddInt32(atomicIndex, 1)); i < len(rs.clusters); i = int(atomic.AddInt32(atomicIndex, 1)) {
				rs.clusters[i].Notify(index)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	//}

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

func (rs *RootCluster) Results() []*CrownCluster {
	var crownClusters = make([]*CrownCluster, 0)
	for _, stemCluster := range rs.clusters {
		for _, crownCluster := range stemCluster.clusters {
			crownClusters = append(crownClusters, crownCluster)
		}
	}
	sort.Slice(crownClusters, func(i,j int) bool {
		return len(crownClusters[i].indexes) > len(crownClusters[j].indexes)
	})
	return crownClusters
}



func (rs *RootCluster) String() string {
	res := ""
	for _, crownCluster := range rs.Results() {
		res += fmt.Sprintln(crownCluster)
	}
	return res
}
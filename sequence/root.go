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
)

type RootCluster struct {
	limit int

	clusters []*StemCluster
	nodeIDToCluster []*StemCluster

	arena *arena.Arena
	strictComparator comparator.Comparator
	elementComparator comparator.Comparator

	notify chan [2]int
}

func NewRootCluster() *RootCluster {
	a := arena.NewArena()
	return &RootCluster{
		limit: 99999999,
		arena: a,
		strictComparator: comparator.NewStrictComparator(a),
		elementComparator: comparator.NewElementComparator(a),
		notify: make(chan [2]int),
	}
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

func (rs *RootCluster) Batch() {
	Init(rs.arena)
	rs.nodeIDToCluster = make([]*StemCluster, len(rs.arena.List))
	rs.consumeNotifications()
	for i := range rs.arena.List {
		rs.Add(i)
	}
	rs.notifyAll()
}

func (rs *RootCluster) consumeNotifications() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(){
			for pair := range rs.notify {
				c := rs.clusters[pair[0]]
				index := pair[1]

				c.m.Lock()
				if len(c.endings) > 0 {
					if index > c.endings[0].last {
						c.addWithCrown(c.endings[0].i, c.endings[0].index)
						c.endings = c.endings[1:]
					}
				}
				c.m.Unlock()
			}
		}()
	}
}

func (rs *RootCluster) notifyAll() {
	for i := range rs.clusters {
		rs.notify <- [2]int{i, len(rs.arena.List)}
	}
}

//func (rs *RootCluster) Rate(index int) float32 { return 0 }
func (rs *RootCluster) Add(index int) bool {
	for i := range rs.clusters {
		rs.notify <- [2]int{i, index}
	}

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
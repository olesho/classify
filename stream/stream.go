package stream

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/olesho/classify/arena"
	"github.com/olesho/classify/comparator"
	"golang.org/x/net/html"
)

// Engine is stream processing entity
type Engine struct {
	Arena        *arena.Arena
	ItemClusters []*ClusterMatrix
	Storage      []*ClusterMatrix

	cmp   comparator.Comparator
	timer *Timer
}

// NewEngine creates new stream processing entity
func NewEngine() *Engine {
	a := arena.NewArena()
	return &Engine{
		Arena: a,

		cmp:   comparator.NewDefaultComparator(a),
		timer: NewTimer(),
	}
}

// LoadFile appends HTML file content
func (e *Engine) LoadFile(fileName string) error {
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
	e.Arena.Append(*n)
	return nil
}

// LoadString appends HTML string
func (e *Engine) LoadString(s string) error {
	n, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return err
	}
	e.Arena.Append(*n)
	return nil
}

// Load appends HTML from reader
func (e *Engine) Load(r io.Reader) error {
	n, err := html.Parse(r)
	if err != nil {
		return err
	}
	e.Arena.Append(*n)
	return nil
}

// Run starts clusterization process
func (e *Engine) Run(windowLength, numCPU int) *Matrix {
	comparator.Init(e.Arena)
	arenaLength := len(e.Arena.List)
	if windowLength <= 0 {
		windowLength = arenaLength
	}
	e.ItemClusters = make([]*ClusterMatrix, arenaLength)

	e.timer.Start()
	fmt.Printf("arena length: %v\n", len(e.Arena.List))

	tempVals := make([][]float32, arenaLength)
	tempIndexes := make([][]int, arenaLength)

	index := new(int32)
	*index = -1
	wg := sync.WaitGroup{}
	wg.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go func() {
			idx := int(atomic.AddInt32(index, 1))
			for idx < len(e.Arena.List) {
				for j := idx + 1; j < arenaLength; j++ {
					val := e.cmp.Cmp(e.Arena.List[idx], e.Arena.List[j])
					if val > 0 {
						tempIndexes[idx] = append(tempIndexes[idx], j)
						tempVals[idx] = append(tempVals[idx], val)
					}
				}
				idx = int(atomic.AddInt32(index, 1))
			}
			wg.Done()
		}()
	}
	wg.Wait()

	alreadyUsed := make([]bool, arenaLength)
	for idx := 0; idx < len(e.Arena.List); idx++ {
		if !alreadyUsed[idx] {
			if len(tempIndexes[idx]) > 0 {
				newCluster := &ClusterMatrix{
					windowSize: windowLength,
				}
				e.ItemClusters[idx] = newCluster
				newCluster.Indexes = append([]int{idx}, tempIndexes[idx]...)
				newCluster.Values = make([][]float32, len(newCluster.Indexes))
				for i, nextIndex := range newCluster.Indexes {
					alreadyUsed[nextIndex] = true
					newCluster.Values[i] = tempVals[nextIndex]
					e.ItemClusters[i] = newCluster
				}
				e.Storage = append(e.Storage, newCluster)
			} else {
				alreadyUsed[idx] = true
			}
		}
	}

	clusters := []*Cluster{}
	for _, matrix := range e.Storage {
		clusters = append(clusters, matrix.Clusters()...)
	}

	tables := make([]Table, len(clusters))
	for i, cluster := range clusters {
		tables[i] = cluster.toTable(e.Arena)
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Rate > clusters[j].Rate
	})

	clusterGroups := groupClusters(e.Arena, tables)
	sort.Slice(clusterGroups, func(i, j int) bool {
		return clusterGroups[i].GroupVolume > clusterGroups[j].GroupVolume
	})

	// transpose
	rm := &Matrix{
		Matrix: make([]Series, len(clusterGroups)),
	}
	for i, g := range clusterGroups {
		rm.Matrix[i] = Series{
			Matrix: transpose(g),
			Arena:  e.Arena,
			Group:  g,
		}
	}

	// trace
	e.timer.Check("all clusters sorted")
	return rm
}

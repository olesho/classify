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

type EngineOpts struct {
	NumCPU int
	WindowSize int
	Comparator comparator.Comparator
}

// Engine is stream processing entity
type Engine struct {
	Arena        *arena.Arena

	windowSize int
	numCPU int
	cmp   comparator.Comparator
	timer *Timer
}

// NewEngine creates new stream processing entity
func NewEngine(opts *EngineOpts) *Engine {
	a := arena.NewArena()
	e := &Engine{
		Arena: a,

		windowSize: opts.WindowSize,
		numCPU: opts.NumCPU,
		cmp: opts.Comparator,
		timer: NewTimer(),
	}
	if e.cmp == nil {
		e.cmp = comparator.NewDefaultComparator(a)
	}
	return e
}

func (e *Engine) SetWindowSize(s int) {
	e.windowSize = s
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

func (e *Engine) genLinesAsync(arenaLength, windowLength, numCPU int) (tempValues [][]float32, tempIndexes [][]int) {
	tempValues = make([][]float32, arenaLength)
	tempIndexes = make([][]int, arenaLength)

	index := new(int32)
	*index = -1
	wg := sync.WaitGroup{}
	wg.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go func() {
			idx := int(atomic.AddInt32(index, 1))
			for idx < len(e.Arena.List) {
				lastIdx := idx + windowLength
				if lastIdx > arenaLength {
					lastIdx = arenaLength
				}
				for j := idx + 1; j < lastIdx; j++ {
					val := e.cmp.Cmp(e.Arena.List[idx], e.Arena.List[j])
					if val > 0 {
						tempIndexes[idx] = append(tempIndexes[idx], j)
						tempValues[idx] = append(tempValues[idx], val)
					}
				}
				idx = int(atomic.AddInt32(index, 1))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return
}

func fillClusterMatrix(cluster *ClusterMatrix, idx int, tempValues [][]float32, tempIndexes [][]int, alreadyUsed []bool) {
	cluster.Indexes = append(cluster.Indexes, idx)
	cluster.Values = append(cluster.Values, tempValues[idx])
	alreadyUsed[idx] = true
	if len(tempIndexes[idx]) > 0 {
		fillClusterMatrix(cluster, tempIndexes[idx][0], tempValues, tempIndexes, alreadyUsed)
	}
}

func mergeLinesInWindow(arenaLength int, tempValues [][]float32, tempIndexes [][]int) []*ClusterMatrix {
	var results []*ClusterMatrix
	alreadyUsed := make([]bool, arenaLength)
	for idx := 0; idx < arenaLength; idx++ {
		if !alreadyUsed[idx] {
			if len(tempIndexes[idx]) > 0 {
				cluster := &ClusterMatrix{
					windowLength: arenaLength,
				}

				cluster.Indexes = []int{idx}
				cluster.Values = make([][]float32, 1)
				cluster.Values[0] = tempValues[idx]
				fillClusterMatrix(cluster, tempIndexes[idx][0], tempValues, tempIndexes, alreadyUsed)
				results = append(results, cluster)
			}
			alreadyUsed[idx] = true
		}
	}
	return results
}

func (e *Engine) mergeClusterMatrix(dest, src *ClusterMatrix) {
	// TODO: should limit dest.Indexes to [:e.WindowSize] to decrease computational complexity
	for i, idx1 := range dest.Indexes {
		node1 := e.Arena.List[idx1]
		values := make([]float32, len(src.Indexes))
		for j, idx2 := range src.Indexes {
			values[j] = e.cmp.Cmp(node1, e.Arena.List[idx2])
		}
		dest.Values[i] = append(dest.Values[i], values...)
	}
	dest.Indexes = append(dest.Indexes, src.Indexes...)
	dest.Values = append(dest.Values, src.Values...)
}

func (e *Engine) mergeClusterMatricesAsync(input []*ClusterMatrix) []*ClusterMatrix {
	pairs := make(chan [2]int)
	locks :=  make([]sync.Mutex, len(input))
	wg := sync.WaitGroup{}
	for n := 0; n < e.numCPU; n++ {
		wg.Add(1)
		go func() {
			for next := range pairs {
				i, j := next[0], next[1]
				locks[i].Lock()
				locks[j].Lock()
				e.mergeClusterMatrix(input[next[0]], input[next[1]])
				input[j] = nil
				locks[i].Unlock()
				locks[j].Unlock()
			}
			wg.Done()
		}()
	}

	for i := 0; i < len(input); i++ {
		if input[i] != nil {
			for j := i + 1; j < len(input); j++ {
				if input[j] != nil {
					val := e.cmp.Cmp(e.Arena.Get(input[i].Indexes[0]), e.Arena.Get(input[j].Indexes[0]))
					if val > 0 {
						pairs <- [2]int{i, j}
					}
				}
			}
		}
	}
	close(pairs)
	wg.Wait()

	refined := make([]*ClusterMatrix, 0)
	for _, m := range input {
		if m != nil {
			refined = append(refined, m)
		}
	}
	return refined
}


func (e *Engine) mergeClusterMatrices(input []*ClusterMatrix) []*ClusterMatrix {
	for i := 0; i < len(input); i++ {
		if input[i] != nil {
			for j := i + 1; j < len(input); j++ {
				if input[j] != nil {
					val := e.cmp.Cmp(e.Arena.Get(input[i].Indexes[0]), e.Arena.Get(input[j].Indexes[0]))
					if val > 0 {
						e.mergeClusterMatrix(input[i], input[j])
						input[j] = nil
					}
				}
			}
		}
	}
	refined := make([]*ClusterMatrix, 0)
	for _, m := range input {
		if m != nil {
			refined = append(refined, m)
		}
	}
	return refined
}

func mergeLines(arenaLength int, tempValues [][]float32, tempIndexes [][]int) []*ClusterMatrix {
	var results []*ClusterMatrix
	alreadyUsed := make([]bool, arenaLength)
	for idx := 0; idx < arenaLength; idx++ {
		if !alreadyUsed[idx] {
			if len(tempIndexes[idx]) > 0 {
				newCluster := &ClusterMatrix{
					windowLength: arenaLength,
				}
				newCluster.Indexes = append([]int{idx}, tempIndexes[idx]...)
				newCluster.Values = make([][]float32, len(newCluster.Indexes))
				for i, nextIndex := range newCluster.Indexes {
					alreadyUsed[nextIndex] = true
					newCluster.Values[i] = tempValues[nextIndex]
				}
				results = append(results, newCluster)
			} else {
				alreadyUsed[idx] = true
			}
		}
	}
	return results
}

// Run process of putting to clusters
func (e *Engine) Run() *Matrix {
	comparator.Init(e.Arena)
	arenaLength := len(e.Arena.List)
	if e.windowSize <= 0 {
		e.windowSize = arenaLength
	}

	e.timer.Start()
	fmt.Printf("arena length: %v\n", len(e.Arena.List))

	tempValues, tempIndexes := e.genLinesAsync(arenaLength, e.windowSize, e.numCPU)
	e.timer.Check("first run finished")

	var clusterMatrices []*ClusterMatrix
	if e.windowSize < arenaLength {
		clusterMatrices = mergeLinesInWindow(arenaLength, tempValues, tempIndexes)
		e.timer.Check("merged lines")

		clusterMatrices = e.mergeClusterMatricesAsync(clusterMatrices)
		e.timer.Check("merged cluster matrices")
	} else {
		clusterMatrices = mergeLines(arenaLength, tempValues, tempIndexes)
		e.timer.Check("merged lines")
	}

	var clusters []*Cluster
	for _, matrix := range clusterMatrices {
		clusters = append(clusters, matrix.Clusters()...)
	}
	e.timer.Check("clusters generated")

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

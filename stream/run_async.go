package stream

import (
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

func (s *Storage) createMatricesAsync() {
	s.NodeToCluster = make([]*Mtx, len(s.Arena.List))
	index := new(int32)
	*index = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			idx := int(atomic.AddInt32(index, 1))
			for idx < len(s.Arena.List) {
				s.next(idx)
				idx = int(atomic.AddInt32(index, 1))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	s.timer.Check("elements compared")
}

func (s *Storage) compareInMatrixAsync(mtx *Mtx) {
	nn := new(int32)
	*nn = -1
	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(dest *Mtx) {
			n := int(atomic.AddInt32(nn, 1))
			for n < len(dest.Values) {
				idx1 := dest.Indexes[n]
				for m, idx2 := range dest.Indexes[n+1:] {
					childrenVal := s.cmpChildren(idx1, idx2)
					dest.Values[n][m] += childrenVal
				}
				n = int(atomic.AddInt32(nn, 1))
			}
			wg.Done()
		}(mtx)
	}
	wg.Wait()
}

func (s *Storage) compareInMatricesAsync() {
	clusterDuplicates := make([]*Mtx, len(s.Clusters))
	for j, mtx := range s.Clusters {
		mtxClone := mtx.Clone()
		if 2*len(mtxClone.Indexes)*(len(mtxClone.Indexes)-1) > runtime.NumCPU() {
			s.compareInMatrixAsync(mtxClone)
		} else {
			s.compareInMatrix(mtxClone)
		}
		clusterDuplicates[j] = mtxClone
	}
	s.Clusters = clusterDuplicates

	for _, c := range s.Clusters {
		for _, idx := range c.Indexes {
			s.NodeToCluster[idx] = c
		}
	}

	s.timer.Check("elements compared including children")
}

func (s *Storage) generateAllClustersAsync() []*Cluster {
	wg := sync.WaitGroup{}
	index := new(int32)
	clusters := make([]*Cluster, 0)
	*index = -1
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			idx := int(atomic.AddInt32(index, 1))
			for idx < len(s.Clusters) {
				generated := s.Clusters[idx].GenerateClusters()
				clusters  = append(clusters, generated...)
				idx = int(atomic.AddInt32(index, 1))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	s.timer.Check("clusters generated")
	return clusters
}

func (s *Storage) clustersToMatrix(clusters []*Cluster) []Series {
	tables := make([]Table, len(clusters))
	for i, cluster := range clusters {
		tables[i] = s.toTable(cluster)
	}

	clusterGroups := groupClusters(s.Arena, tables)
	sort.Slice(clusterGroups, func(i, j int) bool {
		return clusterGroups[i].GroupVolume > clusterGroups[j].GroupVolume
	})

	// transpose
	rm := make([]Series, len(clusterGroups))
	for i, g := range clusterGroups {
		rm[i] = Series{
			Matrix: transpose(g),
			Arena:  s.Arena,
			Group:  g,
		}
	}

	// trace
	s.timer.Check("all clusters sorted")
	return rm
}

func (s *Storage) RunAsync() []Series {
	s.timer.Start()
	s.createMatricesAsync()
	s.compareInMatricesAsync()
	clusters := s.generateAllClustersAsync()
	return s.clustersToMatrix(clusters)
}


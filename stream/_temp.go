package stream

func (s *Storage) _RunAsync() *Matrix {
	s.timer.Start()

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

	lastCluster := s.NewMtx(5119)
	s.Clusters = append(s.Clusters, lastCluster)
	for idx := range s.Arena.List {
		s._next(idx)
	}

	for _, mtx := range s.Clusters[len(s.Clusters)-1:] {
		nn := new(int32)
		*nn = -1
		wg.Add(runtime.NumCPU())
		for i := 0; i < runtime.NumCPU(); i++ {
			go func(i int) {
				n := int(atomic.AddInt32(nn, 1))
				for n < len(mtx.Values) {
					idx1 := mtx.Indexes[n]
					for m, idx2 := range mtx.Indexes[n+1:] {
						mtx.Values[n][m] += s.cmpChildren(idx1, idx2)
					}
					n = int(atomic.AddInt32(nn, 1))
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
	s.timer.Check("elements compared including children")

	var clusters = s.Clusters[len(s.Clusters)-1]._GenerateClusters(s.Arena)

	tables := make([]Table, len(clusters))
	for i, cluster := range clusters {
		tables[i] = s.toTable(cluster)
	}

	clusterGroups := groupClusters(s.Arena, tables)
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
			Arena:  s.Arena,
			Group:  g,
		}
	}

	// trace
	s.timer.Check("all clusters sorted")
	return rm
}

func (s *Storage) __RunAsync() *Matrix {
	s.timer.Start()

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

	for _, mtx := range s.Clusters {
		nn := new(int32)
		*nn = -1
		wg.Add(runtime.NumCPU())
		for i := 0; i < runtime.NumCPU(); i++ {
			go func(i int) {
				n := int(atomic.AddInt32(nn, 1))
				for n < len(mtx.Values) {
					idx1 := mtx.Indexes[n]
					for m, idx2 := range mtx.Indexes[n+1:] {
						mtx.Values[n][m] += s.cmpChildren(idx1, idx2)
					}
					n = int(atomic.AddInt32(nn, 1))
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	}
	s.timer.Check("elements compared including children")

	var clusters = make([]*Cluster, 0)
	*index = -1
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			idx := int(atomic.AddInt32(index, 1))
			for idx < len(s.Clusters) {
				clusters  = append(clusters, s.Clusters[idx]._GenerateClusters(s.Arena)...)
				idx = int(atomic.AddInt32(index, 1))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	s.timer.Check("clusters generated")

	tables := make([]Table, len(clusters))
	for i, cluster := range clusters {
		tables[i] = s.toTable(cluster)
	}

	clusterGroups := groupClusters(s.Arena, tables)
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
			Arena:  s.Arena,
			Group:  g,
		}
	}

	// trace
	s.timer.Check("all clusters sorted")
	return rm
}

func (s *Storage) _next(idx int) {
	if idx == 5119 {
		return
	}
	classes := s.Arena.Get(idx).Classes()
	if len(classes) > 0 {
		if classes[0] == "catalog-grid__cell" {
			s.TryAdd(len(s.Clusters)-1, idx)
		}
	}
}

func (s *Storage) _consolidate(cluster1, cluster2 *Cluster) *Cluster {
	var minVal float32 = s.Find(cluster1.Indexes[0], cluster2.Indexes[0])
	for _, i1 := range cluster1.Indexes {
		for _, i2 := range cluster2.Indexes {
			if val := s.Find(i1, i2); val < minVal {
				minVal = val
			}
		}
	}
	newVolume := float32(len(cluster1.Indexes) + len(cluster2.Indexes)) * minVal
	vol1 := cluster1.Volume()
	vol2 := cluster2.Volume()
	if vol1 + vol2 < newVolume {
		//if cluster1.Volume() < newVolume || cluster2.Volume() < newVolume {
		return &Cluster{
			Indexes: append(cluster1.Indexes, cluster2.Indexes...),
			Rate: minVal,
		}
	}
	return nil
}

func (s *Storage) _consolidateAll(clusters []*Cluster) []*Cluster {
	for i := range clusters {
		for j := i+1; j < len(clusters); j++ {
			if clusters[i] != nil && clusters[j] != nil {
				if newCluster:= s._consolidate(clusters[i], clusters[j]); newCluster != nil {
					clusters[i] = newCluster
					clusters[j] = nil
				}
			}
		}
	}

	filteredClusters := make([]*Cluster, 0)
	for _, c := range clusters {
		if c != nil {
			filteredClusters = append(filteredClusters, c)
		}
	}
	return filteredClusters
}

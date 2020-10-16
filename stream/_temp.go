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


package stream

func (s *Storage) createMatrices() {
	s.NodeToCluster = make([]*Mtx, len(s.Arena.List))
	for idx := 0; idx < len(s.Arena.List); idx++ {
		s.next(idx)
	}
	s.timer.Check("elements compared")
}

func (s *Storage) compareInMatrix(mtx *Mtx) {
	for n := 0; n < len(mtx.Values); n++ {
		idx1 := mtx.Indexes[n]
		for m, idx2 := range mtx.Indexes[n+1:] {
			childrenVal := s.cmpChildren(idx1, idx2)
			mtx.Values[n][m] += childrenVal
		}
	}
}

func (s *Storage) compareInMatrices() {
	clusterDuplicates := make([]*Mtx, len(s.Clusters))
	for j, mtx := range s.Clusters {
		mtxClone := mtx.Clone()
		s.compareInMatrix(mtxClone)
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

func (s *Storage) generateAllClusters() []*Cluster {
	clusters := make([]*Cluster, 0)
	for idx := 0; idx < len(s.Clusters); idx ++ {
		clusters  = append(clusters, s.Clusters[idx].GenerateClusters()...)
	}
	s.timer.Check("clusters generated")
	return clusters
}

func (s *Storage) Run() []Series {
	s.timer.Start()
	s.createMatrices()
	s.compareInMatrices()
	clusters := s.generateAllClusters()
	return s.clustersToMatrix(clusters)
}



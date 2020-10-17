package stream

import "sort"

func (s *Storage) getElemRate(idx1, idx2 int) float32 {
	if s.NodeToCluster[idx1] == s.NodeToCluster[idx2] {
		v := s.NodeToCluster[idx1].Find(idx1, idx2)
		return v
	}
	return 0
}

func (s *Storage) cmpChildren(idx1, idx2 int) float32 {
	n1, n2 := s.Arena.Get(idx1), s.Arena.Get(idx2)
	size1, size2 := len(n1.Children), len(n2.Children)
	rating := make([]rateItem, size1*size2)
	for i1, idx1 := range n1.Children {
		for i2, idx2 := range n2.Children {
			idx := (i1+1)*(i2+1) - 1
			rc := s.getElemRate(idx1, idx2)
			if rc > 0 {
				cv := s.cmpChildren(idx1, idx2)
				rc += cv
			}
			rating[idx].Coincided = rc
			rating[idx].Index1 = i1
			rating[idx].Index2 = i2
		}
	}

	sort.Slice(rating, func(i, j int) bool {
		return rating[i].Coincided > rating[j].Coincided
	})

	flags1 := make([]bool, size1)
	flags2 := make([]bool, size2)
	count := 0
	smallerSize := 0
	if size1 < size2 {
		smallerSize = size1
	} else {
		smallerSize = size2
	}

	var coincided float32 = 0.
	for _, rate := range rating {
		if !flags1[rate.Index1] && !flags2[rate.Index2] {
			if rate.Coincided == 0 {
				break
			}
			coincided += rate.Coincided
			flags1[rate.Index1] = true
			flags2[rate.Index2] = true
			count++
			if count == smallerSize {
				break
			}
		}
	}

	return coincided
}

type rateItem struct {
	Coincided float32
	Index1    int
	Index2    int
}
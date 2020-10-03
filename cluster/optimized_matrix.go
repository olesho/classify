package cluster

import (
	"sync"
	"sync/atomic"
)

// OptimizedMatrix represents similarity between document elements
type OptimizedMatrix struct {
	// Values is H x W float32 matrix where H - size of arena list (all elements), W - sliding window size (number)
	// Values[X][Y] means values of similarity between arena.List[X] and arena.List[X+Y]
	Values      [][]float32
	MaxForIndex []float32
	Excluded    []bool
	windowSize  int
	wg          sync.WaitGroup
}

func Clone(matrix *OptimizedMatrix) *OptimizedMatrix {
	cp := &OptimizedMatrix{}
	cp.Values = make([][]float32, len(matrix.Values))
	for i := range cp.Values {
		cp.Values[i] = make([]float32, len(matrix.Values[i]))
		copy(cp.Values[i], matrix.Values[i])
	}
	cp.MaxForIndex = make([]float32, len(matrix.MaxForIndex))
	copy(cp.MaxForIndex, matrix.MaxForIndex)

	cp.Excluded = make([]bool, len(matrix.Excluded))
	copy(cp.Excluded, matrix.Excluded)

	cp.windowSize = matrix.windowSize
	return cp
}

// NewOptimizedRateMatrixAsync creates similarity matrix using provided comparation function
func NewOptimizedRateMatrixAsync(length, windowLength, numCPU int, cmp func(i, j int) float32) *OptimizedMatrix {
	rm := &OptimizedMatrix{
		Values:      make([][]float32, length),
		MaxForIndex: make([]float32, length),
		Excluded:    make([]bool, length),

		windowSize: windowLength,
		wg:         sync.WaitGroup{},
	}

	index := new(int32)
	*index = -1
	wg := sync.WaitGroup{}
	wg.Add(4)
	for i := 0; i < numCPU; i++ {
		go func() {
			idx := int(atomic.AddInt32(index, 1))
			for idx < len(rm.Values) {
				currWindowLength := windowLength
				if idx+windowLength >= length {
					currWindowLength = length - idx
				}

				rm.Values[idx] = make([]float32, currWindowLength)
				var max float32 = .0
				for j := range rm.Values[idx : idx+currWindowLength-1] {
					val := cmp(idx, idx+j+1)
					if val > max {
						max = val
					}
					rm.Values[idx][j] = val
				}
				rm.MaxForIndex[idx] = max
				idx = int(atomic.AddInt32(index, 1))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return rm
}

// NewOptimizedRateMatrix creates similarity matrix using provided comparation function
func NewOptimizedRateMatrix(length, windowLength int, cmp func(i, j int) float32) *OptimizedMatrix {
	rm := &OptimizedMatrix{
		Values:      make([][]float32, length),
		MaxForIndex: make([]float32, length),
		Excluded:    make([]bool, length),
		windowSize:  windowLength,
	}
	for i := range rm.Values {
		currWindowLength := windowLength
		if i+windowLength >= length {
			currWindowLength = length - i
		}

		rm.Values[i] = make([]float32, currWindowLength)
		var max float32 = .0
		for j := range rm.Values[i : i+currWindowLength-1] {
			val := cmp(i, i+j+1)
			if val > max {
				max = val
			}
			rm.Values[i][j] = val
		}
		rm.MaxForIndex[i] = max
	}
	return rm
}

// Get returns similarity by given indexes
func (m *OptimizedMatrix) Get(idx1, idx2 int) float32 {
	if idx1 < idx2 {
		diff := idx2 - idx1 - 1
		if diff < m.windowSize {
			return m.Values[idx1][diff]
		}
		return 0
	}
	if idx1 == idx2 {
		return 0
	}
	diff := idx1 - idx2 - 1
	if diff < m.windowSize {
		return m.Values[idx2][diff]
	}
	return 0
}

// Max returns maximum similarity value and indexes
func (m *OptimizedMatrix) Max() (max float32, maxi, maxj int) {
	max = .0
	maxi = -1
	maxj = -1

	var maxForIndex float32 = .0
	for i := range m.MaxForIndex {
		if m.MaxForIndex[i] > maxForIndex {
			maxForIndex = m.MaxForIndex[i]
			maxi = i
		}
	}
	if maxi > -1 {
		for j, cell := range m.Values[maxi] {
			if cell > max {
				max = cell
				maxj = maxi + j + 1
			}
		}
	}
	return
}

func maxInSlice(s []float32) float32 {
	var max float32 = .0
	for _, v := range s {
		if v > max {
			max = v
		}
	}
	return max
}

// Exclude marks index as already used
func (m *OptimizedMatrix) Exclude(idx1, idx2 int) {
	if idx1 < idx2 {
		diff := idx2 - idx1 - 1
		if diff < m.windowSize {
			m.Values[idx1][diff] = 0
		}
		return
	}
	if idx1 == idx2 {
		return
	}
	diff := idx1 - idx2 - 1
	if diff < m.windowSize {
		m.Values[idx2][diff] = 0
	}
	return
}

// Exclude marks index as already used
func (m *OptimizedMatrix) ExcludeRow(index int) {
	m.Excluded[index] = true
	m.MaxForIndex[index] = 0
	for i := range m.Values[:index] {
		diff := index - i - 1
		if m.Values[i][diff] == m.MaxForIndex[i] {
			m.Values[i][diff] = 0
			m.MaxForIndex[i] = maxInSlice(m.Values[i])
		} else {
			m.Values[i][diff] = 0
		}
	}
}

func (m *OptimizedMatrix) IsRowExcluded(index int) bool {
	if index >= len(m.Excluded) {
		return true
	}
	return m.Excluded[index]
}

func (m *OptimizedMatrix) Candidates(idxs []int) (pairIdxs []int) {
	for _, idx := range idxs {
		pairIdxs = append(pairIdxs, m.candidatesForIdx(idx)...)
	}
	return
}

func (m *OptimizedMatrix) candidatesForIdx(idx int) (pairIdxs []int) {
	for pairIdx, val := range m.Values[idx] {
		if val > 0 {
			pairIdxs = append(pairIdxs, idx+pairIdx+1)
		}
	}
	for pairIdx, row := range m.Values {
		shiftIdx := idx - pairIdx - 1
		if shiftIdx >= len(row) || shiftIdx < 0 {
			continue
		}
		val := row[shiftIdx]
		if val > 0 {
			pairIdxs = append(pairIdxs, pairIdx)
		}
	}
	return
}

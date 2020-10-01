package cluster

import "time"

// ComparableList is a list of results of node comparations
type ComparableList interface {
	Cmp(i, j int) float64
	Exclude(index int)
	IsExcluded(index int) bool
	Candidates(idx int) (startingIndex int, values []float64)
}

// RateMatrix is most straghtforward ComparableList implementation
type RateMatrix struct {
	Rows     []RateRow
	Excluded []bool

	excludedCount int
	startedAt     time.Time
}

type RateRow []float64

func NewRateMatrix(size1, size2 int, cmp func(i, j int) float64) *RateMatrix {
	rm := &RateMatrix{
		excludedCount: 0,
		startedAt:     time.Now(),
	}
	cells := make([]RateRow, size1)
	off := make([]bool, size1)
	for i := range cells {
		cells[i] = make([]float64, size2)
		for j := range cells[i] {
			val := cmp(i, j)
			cells[i][j] = val
		}
	}
	rm.Rows = cells
	rm.Excluded = off
	return rm
}

func (m *RateMatrix) Cmp(idx1, idx2 int) float64 {
	if idx1 < idx2 {
		return m.Rows[idx1][idx2]
	}
	return m.Rows[idx2][idx1]
}

func (m *RateMatrix) Max() (max float64, maxi, maxj int) {
	max = .0
	maxi = -1
	maxj = -1
	for i, row := range m.Rows {
		if !m.IsExcluded(i) {
			for j, cell := range row {
				if !m.IsExcluded(j) {
					if cell > max {
						max = cell
						maxi = i
						maxj = j
					}
				}
			}
		}
	}
	return
}

func (m *RateMatrix) Exclude(index int) {
	m.Excluded[index] = true
	m.excludedCount++
}

func (m *RateMatrix) IsExcluded(index int) bool {
	return m.Excluded[index]
}

func (m *RateMatrix) Candidates(idx int) (startingIndex int, values []float64) {
	return 0, m.Rows[idx]
}

package cluster

import "time"

type RateMatrix struct {
	Rows        []RateRow
	RowExcluded []bool
	ColExcluded []bool

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
	offRows := make([]bool, size1)
	offCols := make([]bool, size2)
	for i := range cells {
		cells[i] = make([]float64, size2)
		for j := range cells[i] {
			val := cmp(i, j)
			cells[i][j] = val
		}
	}
	rm.Rows = cells
	rm.RowExcluded = offRows
	rm.ColExcluded = offCols
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
		if !m.RowExcluded[i] {
			for j, cell := range row {
				if !m.ColExcluded[j] {
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

func (m *RateMatrix) ExcludeRows(index int) {
	m.RowExcluded[index] = true
	m.excludedCount++
}

func (m *RateMatrix) ExcludeCols(index int) {
	m.ColExcluded[index] = true
}

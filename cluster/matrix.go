package cluster

import (
	"time"

	"github.com/olesho/classify/arena"
)

// ComparableList is a list of results of node comparations
type ComparableList interface {
	Get(i, j int) float32
	Exclude(i, j int)
	ExcludeRow(index int)
	IsRowExcluded(index int) bool
	Candidates(idxs []int) (pairIdxs []int)
}

type MergerList interface {
	MergeAll(arena *arena.Arena, indexes []int) *arena.Arena
	MergeIntoTemplate(mainArena, templateArena *arena.Arena, mainIdx, templateIdx int)
}

// RateMatrix is most straghtforward ComparableList implementation
type RateMatrix struct {
	Values   []RateRow
	Excluded []bool

	excludedCount int
	startedAt     time.Time
}

type RateRow []float32

func NewRateMatrix(size1, size2 int, cmp func(i, j int) float32) *RateMatrix {
	rm := &RateMatrix{
		excludedCount: 0,
		startedAt:     time.Now(),
	}
	cells := make([]RateRow, size1)
	off := make([]bool, size1)
	for i := range cells {
		cells[i] = make([]float32, size2)
		for j := range cells[i] {
			val := cmp(i, j)
			cells[i][j] = val
		}
	}
	rm.Values = cells
	rm.Excluded = off
	return rm
}

func (m *RateMatrix) Get(idx1, idx2 int) float32 {
	if idx1 < idx2 {
		return m.Values[idx1][idx2]
	}
	return m.Values[idx2][idx1]
}

func (m *RateMatrix) Max() (max float32, maxi, maxj int) {
	max = .0
	maxi = -1
	maxj = -1
	for i, row := range m.Values {
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

func (m *RateMatrix) Candidates(idxs []int) (pairIdxs []int) {
	for _, idx := range idxs {
		pairIdxs = append(pairIdxs, m.candidatesForIdx(idx)...)
	}
	return
}

func (m *RateMatrix) candidatesForIdx(idx int) (pairIdxs []int) {
	for pairIdx, val := range m.Values[idx] {
		if val > 0 {
			pairIdxs = append(pairIdxs, pairIdx)
		}
	}
	for pairIdx, row := range m.Values {
		val := row[idx]
		if val > 0 {
			pairIdxs = append(pairIdxs, pairIdx)
		}
	}
	return
}

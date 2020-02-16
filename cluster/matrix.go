package cluster

type RateMatrix struct {
	Rows []RateRow
	RowFlags []bool
	ColFlags []bool
}

type RateRow []float64

func NewRateMatrix(size1, size2 int, cmp func(i, j int) float64) *RateMatrix {
	cells := make([]RateRow, size1)
	offRows := make([]bool, size1)
	offCols := make([]bool, size2)
	for i  := range cells {
		cells[i] = make([]float64, size2)
		for j := range cells[i] {
			cells[i][j] = cmp(i, j)
		}
	}
	return &RateMatrix{ cells, offRows, offCols }
}

func (m *RateMatrix) Max() (max float64, maxi, maxj int) {
	max = .0
	maxi = -1
	maxj = -1
	for i, row := range m.Rows {
		if !m.RowFlags[i] {
			for j, cell := range row {
				if !m.ColFlags[j] {
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

func (m *RateMatrix) OffRows(index int) {
	m.RowFlags[index] = true
}


func (m *RateMatrix) OffCols(index int) {
	m.ColFlags[index] = true
}

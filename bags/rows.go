package bags

import classify "github.com/olesho/class"

type Row []*classify.Node

type Matrix struct {
	Matrix [][]Row
	Arena  *classify.Arena
}

type Rank struct {
	Matrix []Row
	Arena  *classify.Arena
}

func (m *Matrix) Rank(rank int) *Rank {
	if len(m.Matrix) > rank {
		return &Rank{m.Matrix[rank], m.Arena}
	}
	return nil
}

func (m *Rank) isFieldUniform(index int) bool {
	val := m.Arena.StringifyInformation(m.Matrix[0][index].Id)
	for _, row := range m.Matrix {
		if m.Arena.StringifyInformation(row[index].Id) != val {
			return false
		}
	}
	return true
}

func (m *Rank) Uniform() *Rank {
	result := &Rank{Arena: m.Arena, Matrix: make([]Row, len(m.Matrix))}
	uniformity := make([]bool, len(m.Matrix[0]))
	for i := 0; i < len(m.Matrix[0]); i++ {
		uniformity[i] = m.isFieldUniform(i)
	}

	for rowIndex, row := range m.Matrix {
		for i, isUniform := range uniformity {
			if isUniform {
				result.Matrix[rowIndex] = append(result.Matrix[rowIndex], row[i])
			}
		}
	}

	return result
}

func (m *Rank) Nonuniform() *Rank {
	result := &Rank{Arena: m.Arena, Matrix: make([]Row, len(m.Matrix))}
	uniformity := make([]bool, len(m.Matrix[0]))
	for i := 0; i < len(m.Matrix[0]); i++ {
		uniformity[i] = m.isFieldUniform(i)
	}

	for rowIndex, row := range m.Matrix {
		for i, isUniform := range uniformity {
			if !isUniform {
				result.Matrix[rowIndex] = append(result.Matrix[rowIndex], row[i])
			}
		}
	}

	return result
}

// type Pattern struct {
// 	arena *classify.Arena
// 	row   Row
// }

// func clone(a *classify.Arena, row Row) (*classify.Arena, Row) {
// 	resArena := &classify.Arena{}
// 	for _, item := range row {
// 		resArena = a.Clone(item.Id)
// 	}
// }

// func (p *Pattern) Merge(a *classify.Arena, r Row) {
// 	if len(p.row) == 0 {
// 		p.arena, p.row = clone(a, r)
// 		return
// 	}
// }

// func UniformPattern(input []Row) Pattern {
// 	rowSize := len(input)
// 	for i := 0; i < rowSize; i++ {
// 		for _, row := range input {

// 		}
// 	}
// }

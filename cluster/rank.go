package cluster

import (
	"github.com/olesho/classify"
)

type Row []*classify.Node

type Matrix struct {
	Matrix [][]Row
	Arena  *classify.Arena
}

type Series struct {
	Matrix []Row
	Arena  *classify.Arena
}

func (m *Matrix) Nth(rank int) *Series {
	if len(m.Matrix) > rank {
		return &Series{m.Matrix[rank], m.Arena}
	}
	return nil
}

func (m *Series) isFieldUniform(index int) bool {
	val := m.Arena.StringifyInformation(m.Matrix[0][index].Id)
	for _, row := range m.Matrix {
		if m.Arena.StringifyInformation(row[index].Id) != val {
			return false
		}
	}
	return true
}

func (m *Series) Uniform() *Series {
	result := &Series{Arena: m.Arena, Matrix: make([]Row, len(m.Matrix))}
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

func (m *Series) Nonuniform() *Series {
	result := &Series{Arena: m.Arena, Matrix: make([]Row, len(m.Matrix))}
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

func (s *Series) Informative() *Series {
	result := &Series{Arena: s.Arena, Matrix: make([]Row, len(s.Matrix))}
	for rowIndex, row := range s.Matrix {
		//for _, item := range row {
		//	fmt.Println(item.Data)
		//}
		result.Matrix[rowIndex] = row
	}
	return result
}

func (s *Series) String() string {
	stopper := 0
	//for i, n := range s.Arena.List {
	//	if n.Data == "table" && n.HasClass("itemlist") {
	//		stopper = i
	//	}
	//}

	result := ""
	for _, row := range s.Matrix {
		for _, item := range row {
			chain := s.Arena.Chain(item.Id, stopper)
			result += chain.XPath() + "\n"
		}
		result += "--------------------------------------------------\n"
	}
	return result
}

func transpose(group *BagGroup) []Row {
	size := len(group.Bags[0].Members)
	newGroup := make([]Row, size)
	for i := 0; i < size; i++ {
		row := Row{}
		for _, bag := range group.Bags {
			row = append(row, bag.Members[i])
		}
		newGroup[i] = row
	}
	return newGroup
}

func (s *Series) Patterns() *classify.Template {
	groupSize := len(s.Matrix[0])
	for _, row := range s.Matrix[1:] {
		if groupSize != len(row) {
			panic("Rows size not equal")
			return nil
		}
	}

	var templates = make([]classify.Template, len(s.Matrix))
	for templIdx, row := range s.Matrix {
		templates[templIdx] = classify.Template{Chains: make([]classify.Chain, len(row))}
		for i, n := range row {
			templates[templIdx].Chains[i] = s.Arena.Chain(n.Id, 0)
		}
	}

	template := classify.MergeTemplates(&templates[0], &templates[1])
	for _, next := range templates[2:] {
		template = classify.MergeTemplates(template, &next)
	}
	return template
}
package stream

import (
	"strings"

	"github.com/olesho/classify/arena"
)

type Nodes []*arena.Node

//type Matrix struct {
//	Matrix []Series
//}

type Series struct {
	Matrix []Nodes
	Arena  *arena.Arena
	Group  *ClusterGroup
	//Volume float32
	//WholesomeVolume float32
	//Size int
}

func (m *Series) isFieldUniform(index int) bool {
	val := strings.Join(m.Arena.StringifyInformation(m.Matrix[0][index].Id), " ")
	for _, row := range m.Matrix {
		if strings.Join(m.Arena.StringifyInformation(row[index].Id), " ") != val {
			return false
		}
	}
	return true
}

func (m *Series) Uniform() *Series {
	result := &Series{Arena: m.Arena, Matrix: make([]Nodes, len(m.Matrix))}
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
	result := &Series{Arena: m.Arena, Matrix: make([]Nodes, len(m.Matrix))}
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

func (s *Series) String() string {
	stopper := 0
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

func transpose(group *ClusterGroup) []Nodes {
	size := len(group.Clusters[0].Members)
	newGroup := make([]Nodes, size)
	for i := 0; i < size; i++ {
		row := Nodes{}
		for _, bag := range group.Clusters {
			row = append(row, bag.Members[i])
		}
		newGroup[i] = row
	}
	return newGroup
}

func (s *Series) Patterns() *arena.Template {
	groupSize := len(s.Matrix[0])
	for _, row := range s.Matrix[1:] {
		if groupSize != len(row) {
			panic("Rows size not equal")
			return nil
		}
	}

	var templates = make([]arena.Template, len(s.Matrix))
	for templIdx, row := range s.Matrix {
		templates[templIdx] = arena.Template{Chains: make([]arena.Chain, len(row))}
		for i, n := range row {
			templates[templIdx].Chains[i] = s.Arena.Chain(n.Id, 0)
		}
	}

	template := arena.MergeTemplates(&templates[0], &templates[1])
	for _, next := range templates[2:] {
		template = arena.MergeTemplates(template, &next)
	}
	return template
}

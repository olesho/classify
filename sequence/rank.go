package sequence

import (
	"strings"

	"github.com/olesho/classify/arena"
)

type Nodes []*arena.Node

type Series struct {
	//TransposedFields []FieldSet
	TransposedValues [][]string
	TransposedNodes  []Nodes
	Arena            *arena.Arena
	Group            *ClusterGroup
	//Volume float32
	//WholesomeVolume float32
	//Size int
}

func (m *Series) isFieldUniform(index int) bool {
	val := strings.Join(m.Arena.StringifyInformation(m.TransposedNodes[0][index].Id), " ")
	for _, row := range m.TransposedNodes {
		if strings.Join(m.Arena.StringifyInformation(row[index].Id), " ") != val {
			return false
		}
	}
	return true
}

func (m *Series) Uniform() *Series {
	result := &Series{Arena: m.Arena, TransposedNodes: make([]Nodes, len(m.TransposedNodes))}
	uniformity := make([]bool, len(m.TransposedNodes[0]))
	for i := 0; i < len(m.TransposedNodes[0]); i++ {
		uniformity[i] = m.isFieldUniform(i)
	}

	for rowIndex, row := range m.TransposedNodes {
		for i, isUniform := range uniformity {
			if isUniform {
				result.TransposedNodes[rowIndex] = append(result.TransposedNodes[rowIndex], row[i])
			}
		}
	}

	return result
}

func (m *Series) Nonuniform() *Series {
	result := &Series{Arena: m.Arena, TransposedNodes: make([]Nodes, len(m.TransposedNodes))}
	uniformity := make([]bool, len(m.TransposedNodes[0]))
	for i := 0; i < len(m.TransposedNodes[0]); i++ {
		uniformity[i] = m.isFieldUniform(i)
	}

	for rowIndex, row := range m.TransposedNodes {
		for i, isUniform := range uniformity {
			if !isUniform {
				result.TransposedNodes[rowIndex] = append(result.TransposedNodes[rowIndex], row[i])
			}
		}
	}

	return result
}

func (s *Series) String() string {
	stopper := 0
	result := ""
	for _, row := range s.TransposedNodes {
		for _, item := range row {
			chain := s.Arena.Chain(item.Id, stopper)
			result += chain.XPath() + "\n"
		}
		result += "--------------------------------------------------\n"
	}
	return result
}

func transpose(group *ClusterGroup) *Series {
	size := len(group.Clusters[0].Members)
	transposedNodes := make([]Nodes, size)
	for i := 0; i < size; i++ {
		row := Nodes{}
		for _, bag := range group.Clusters {
			row = append(row, bag.Members[i])
		}
		transposedNodes[i] = row
	}

	//transposedFields := make([][]Field, size)
	transposedValues := make([][]string, size)
	for i := 0; i < size; i++ {
		for _, cluster := range group.Clusters {
			for _, field := range cluster.FieldSets {
				//transposedFields[i] = append(transposedFields[i], Field{field.Type, field.Content[i]})
				transposedValues[i] = append(transposedValues[i], field.Content[i])
			}
		}
	}

	return &Series{
		Group: group,
		//TransposedFields: transposedFields,
		TransposedValues: transposedValues,
		TransposedNodes:  transposedNodes,
	}
}

func (s *Series) Patterns() *arena.Template {
	groupSize := len(s.TransposedNodes[0])
	for _, row := range s.TransposedNodes[1:] {
		if groupSize != len(row) {
			panic("Rows size not equal")
			return nil
		}
	}

	var templates = make([]arena.Template, len(s.TransposedNodes))
	for templIdx, row := range s.TransposedNodes {
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

func (s *Series) XPaths() []string {
	var results = make([]string, 0)
	for _, cluster := range s.Group.Clusters {
		for _, field := range cluster.FieldSets {
			chains := make([]arena.Chain, 0)
			for _, id := range field.IDs {
				chains = append(chains, cluster.Arena.Chain(id, 0))
			}
			resultChain := chains[0]
			for _, chain := range chains[1:] {
				resultChain = arena.MergeChains(chain, resultChain)
			}
			results = append(results, resultChain.XPath())
		}
	}
	return results
}

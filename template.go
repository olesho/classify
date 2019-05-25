// classify project classify.go
package classify

import (
	"golang.org/x/net/html"
)

type Template struct {
	Chains []Chain
}
type Chain []*Node

func max(m [][]float64, off_i, off_j []bool) (i, j int) {
	max := float64(0)
	i, j = -1, -1
	for ii, row := range m {
		for jj, _ := range row {
			if !off_i[ii] && !off_j[jj] {
				if row[jj] > max {
					max = m[ii][jj]
					i = ii
					j = jj
				}
			}
		}
	}
	return
}

func MergeTemplates(c1, c2 *Template) *Template {
	m := make([][]float64, len(c1.Chains))
	for i, _ := range m {
		m[i] = make([]float64, len(c2.Chains))
		for j, _ := range m[i] {
			m[i][j] = cmpChains(c1.Chains[i], c2.Chains[j])
		}
	}

	result := &Template{
		Chains: make([]Chain, 0),
	}

	off_i := make([]bool, len(c1.Chains))
	off_j := make([]bool, len(c2.Chains))
	max_i, max_j := max(m, off_i, off_j)
	for max_i > -1 && max_j > -1 {
		off_i[max_i] = true
		off_j[max_j] = true
		result.Chains = append(result.Chains, mergeChains(c1.Chains[max_i], c2.Chains[max_j]))
		max_i, max_j = max(m, off_i, off_j)
	}

	return result
}

func NewTemplate(a *Arena, nId int) *Template {
	return &Template{
		Chains: a.Chains(nId),
	}
}

func (a *Arena) Chains(nId int) []Chain {
	infoList := a.Infomative(nId)
	chains := make([]Chain, len(infoList))
	for i, infoId := range infoList {
		chains[i] = a.Chain(infoId, nId)
	}
	return chains
}

func (a *Arena) Infomative(nId int) []int {
	var r []int
	if a.List[nId].isInformative() {
		r = append(r, nId)
		return r
	}

	for _, id := range a.List[nId].Children {
		r = append(r, a.Infomative(id)...)
	}

	return r
}

func (a *Arena) Chain(nId, stopper int) Chain {
	c := make([]*Node, 0)
	return a.chain(c, nId, stopper)
}

func (a *Arena) chain(ch []*Node, nId, stopper int) []*Node {
	n := a.Get(nId)
	if n.Parent > 0 && n.Parent != stopper {
		return a.chain(append(ch, n), n.Parent, stopper)
	}
	return append(ch, n)
}

func cmpChains(c1, c2 Chain) float64 {
	min_len := len(c1)
	if len(c2) < min_len {
		min_len = len(c2)
	}

	r := CmpResult{}
	for i := 0; i < min_len; i++ {
		res := CmpShallow(c1[i], c2[i])
		if res == nil {
			break
		}
		if res.Count == 0 {
			break
		}
		r.Append(*res)
	}

	return r.Rate()
}

func mergeChains(c1, c2 Chain) Chain {
	min_len := len(c1)
	if len(c2) < min_len {
		min_len = len(c2)
	}

	res := make([]*Node, 0)
	for i := 0; i < min_len; i++ {
		next := MergeShallow(c1[i], c2[i])
		if next.Data == "" && next.Type == html.ElementNode {
			break
		}
		res = append(res, next)
	}

	return res
}

func (a *Arena) XPath(n int, stopper int) string {
	parent := a.Get(n).Parent
	if parent > 0 && parent != stopper {
		return a.XPath(parent, stopper) + a.Get(n).String()
	}
	return a.Get(n).String()
}

func (a *Arena) PathArray(n int) []int {
	init := make([]int, 0)
	return a.pathArray(init, n)
}

// iterate all nodes up to root
func (a *Arena) pathArray(init []int, n int) []int {
	parent := a.Get(n).Parent
	if parent > 0 {
		return a.pathArray(append(init, n), parent)
	}
	return append(init, n)
}

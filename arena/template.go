// classify project classify.go
package arena

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Template struct {
	Chains []Chain
}
type Chain []*Node

func max(m [][]float32, off_i, off_j []bool) (i, j int) {
	max := float32(0)
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
	m := make([][]float32, len(c1.Chains))
	for i, _ := range m {
		m[i] = make([]float32, len(c2.Chains))
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
		result.Chains = append(result.Chains, MergeChains(c1.Chains[max_i], c2.Chains[max_j]))
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
	infoList := a.Wholesome(nId)
	chains := make([]Chain, len(infoList))
	for i, infoId := range infoList {
		chains[i] = a.Chain(infoId, nId)
	}
	return chains
}

func (a *Arena) Wholesome(nId int) []int {
	var r []int
	if a.List[nId].isWholesome() {
		r = append(r, nId)
		return r
	}

	for _, id := range a.List[nId].Children {
		r = append(r, a.Wholesome(id)...)
	}

	return r
}

func (a *Arena) Chain(nId, stopper int) Chain {
	c := make([]*Node, 0)
	return a.chain(c, nId, stopper)
}

//func (a *Arena) ChainIDXs(nId, stopper int) []int {
//	if nId == stopper {
//		return nil
//	}
//	c := []int{nId}
//	for n := a.Get(nId); n.Parent != stopper; n = a.Get(n.Parent) {
//		c = append(c, n.Parent)
//	}
//	return c
//}

func (a *Arena) chain(ch []*Node, nId, stopper int) []*Node {
	if nId == stopper {
		return nil
	}
	n := a.Get(nId)
	if n.Parent > 0 && n.Parent != stopper {
		return a.chain(append(ch, n), n.Parent, stopper)
	}
	return append(ch, n)
}

func cmpChains(c1, c2 Chain) float32 {
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

func MergeChains(c1, c2 Chain) Chain {
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

func MergeShallow(n1 *Node, n2 *Node) *Node {
	if n1.Type == n2.Type {
		r := Node{
			Type: n1.Type,
			Attr: make([]html.Attribute, 0),
		}

		if n1.Data == n2.Data {
			r.Data = n1.Data
			for _, a1 := range n1.Attr {
				for _, a2 := range n2.Attr {
					if a1.Key == a2.Key && a1.Key != "class" {
						attr := html.Attribute{Key: a1.Key}
						if a1.Val == a2.Val {
							attr.Val = a1.Val
						}
						r.Attr = append(r.Attr, attr)
					}
				}
			}

			classes1 := n1.Classes()
			classes2 := n2.Classes()
			for _, c1 := range classes1 {
				for _, c2 := range classes2 {
					if c1 == c2 {
						r.AddClass(c1)
					}
				}
			}

		}

		/*
			if n1.DataArray != nil && n2.DataArray != nil {
				r.DataArray = append(n1.DataArray, n2.DataArray...)
			}
		*/

		return &r
	}
	return nil
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

func (c Chain) XPath() string {
	sections := []string{}
	for i := len(c) - 1; i >= 0; i-- {
		sections = append(sections, elemToXPath(c[i]))
	}
	return strings.Join(sections, "/")
}

func elemToXPath(el *Node) string {
	if el.Type == html.ElementNode {
		return el.Data + attrsToXPath(el.Attr)
	}
	if el.Type == html.TextNode {
		return fmt.Sprintf("@text=%s", el.Data)
	}
	return ""
}

func attrsToXPath(attrs []html.Attribute) string {
	sections := []string{}
	for _, attr := range attrs {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			for _, class := range classes {
				sections = append(sections, fmt.Sprintf(`contains(@class, "%s")`, class))
			}
			continue
		}
		if attr.Key != "" && attr.Val != "" {
			sections = append(sections, fmt.Sprintf(`@%s="%s"`, attr.Key, attr.Val))
		}
	}
	if len(sections) > 0 {
		return "[" + strings.Join(sections, " and ") + "]"
	}
	return ""
}

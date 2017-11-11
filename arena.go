// classify project classify.go
package classify

import (
	//	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Arena struct {
	List []Node
}

func (a *Arena) Get(id int) Node {
	return a.List[id]
}

func (a *Arena) FindByAttr(k string, v string) *Node {
	for _, n := range a.List {
		for _, attr := range n.Attr {
			if attr.Key == k && attr.Val == v {
				return &n
			}
		}
	}
	return nil
}

func (a *Arena) FindNodeIdByAttr(k string, v string) int {
	for id, n := range a.List {
		for _, attr := range n.Attr {
			if attr.Key == k && attr.Val == v {
				return id
			}
		}
	}
	return -1
}

func (a *Arena) AddChild(p int, c int) {
	a.List[p].Children = append(a.List[p].Children, c)
	a.List[c].Parent = p
}

func NewArena(root html.Node) *Arena {
	result := NewArenaRoot()
	result.transform(0, root)
	return result
}

func NewArenaRoot() *Arena {
	return &Arena{
		List: make([]Node, 0),
	}
}

func (a *Arena) HasParent(child, parent int) bool {
	n := a.Get(child)
	for n.Parent != 0 {
		if n.Parent == parent {
			return true
		}
		n = a.Get(n.Parent)
	}
	return false
}

func (a *Arena) transform(node_index int, n html.Node) {
	if n.Type == html.CommentNode ||
		n.Type == html.ErrorNode ||
		(n.Type == html.ElementNode && strings.ToLower(n.Data) == "noscript") ||
		(n.Type == html.ElementNode && strings.ToLower(n.Data) == "script") ||
		(n.Type == html.TextNode && strings.TrimSpace(n.Data) == "") {
		return
	}

	a.List = append(a.List, *NewNode(n))
	currentId := len(a.List) - 1
	a.AddChild(node_index, currentId)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.transform(currentId, *c)
	}
}

func clone(src, dst *Arena, srcId, dstId int, dstParentId int) {
	dst.List[dstId].Type = src.List[srcId].Type
	dst.List[dstId].Data = src.List[srcId].Data
	dst.List[dstId].Attr = src.List[srcId].Attr
	dst.AddChild(dstParentId, dstId)
	for _, c := range src.List[srcId].Children {
		dst.List = append(dst.List, Node{})
		dstChildId := len(dst.List) - 1
		clone(src, dst, c, dstChildId, dstId)
		//dst.AddChild(dstParentId, dstChildId)
	}
}

func (a *Arena) Clone(srcId int) (dest *Arena) {
	dest = &Arena{make([]Node, 2)}
	clone(a, dest, srcId, 1, 0)
	return dest
}

// clone Node; dstId should aready exist in 'dest'
/*
func (a *Arena) Clone(srcId int, dstId int, dest *Arena) {
	for _, c := range a.List[srcId].Children {
		dest.list = append(dest.list, Node{
			Type: a.List[c].Type,
			Data: a.List[c].Data,
			Attr: a.List[c].Attr,
		})
		destChildId := len(dest.list) - 1
		dest.AddChild(dstId, destChildId)
		a.Clone(c, destChildId, dest)
	}
}
*/

func (a *Arena) PrintList() string {
	res := ""
	for _, n := range a.List {
		if n.Type == html.TextNode {
			res += "text:" + strings.TrimSpace(n.Data) + "\n"
		} else {
			res += n.Data + ":" + n.printAttr() + "\n"
		}
	}
	return res
}

func (a *Arena) CmpColumn(n1 int, n2 int) *CmpResult {
	var cnt float64
	curr_el1 := n1
	curr_el2 := n2
	p1 := a.Get(curr_el1).Parent
	p2 := a.Get(curr_el2).Parent
	next_rate := CmpShallow(a.Get(curr_el1), a.Get(curr_el2))

	if next_rate == nil {
		return next_rate
	}

	if p1 == 0 || p2 == 0 {
		return next_rate
	}

	cnt++

	var total CmpResult
	for p1 != p2 {
		total.Append(*next_rate)

		curr_el1 = a.Get(curr_el1).Parent
		curr_el2 = a.Get(curr_el2).Parent
		p1 = a.Get(curr_el1).Parent
		p2 = a.Get(curr_el2).Parent
		next_rate = CmpShallow(a.Get(curr_el1), a.Get(curr_el2))

		cnt++

		if next_rate == nil {
			return next_rate
		}

		if p1 == 0 || p2 == 0 {
			return next_rate
		}
	}

	total.Append(*next_rate)
	return &total
}

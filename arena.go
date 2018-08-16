// classify project classify.go
package classify

import (
	"strings"

	"golang.org/x/net/html"
)

type Arena struct {
	List []Node
}

func (a *Arena) Get(id int) Node {
	return a.List[id]
}

func (a *Arena) NodesByClass(className string) []Node {
	res := []Node{}
	for _, n := range a.List {
		if n.HasClass(className) {
			res = append(res, n)
		}
	}
	return res
}

func (a *Arena) IndexesByClass(className string) []int {
	res := []int{}
	for i, n := range a.List {
		if n.HasClass(className) {
			res = append(res, i)
		}
	}
	return res
}

func (a *Arena) NodesByAttr(k, v string) []Node {
	res := []Node{}
	for _, n := range a.List {
		for _, attr := range n.Attr {
			if attr.Key == k && attr.Val == v {
				res = append(res, n)
			}
		}
	}
	return res
}

func (a *Arena) IndexesByAttr(k string, v string) []int {
	res := []int{}
	for id, n := range a.List {
		for _, attr := range n.Attr {
			if attr.Key == k && attr.Val == v {
				res = append(res, id)
			}
		}
	}
	return res
}

func (a *Arena) AddChild(p int, c int) {
	a.List[p].Children = append(a.List[p].Children, c)
	a.List[c].Parent = p
}

func NewArenaHtml(data string) (*Arena, error) {
	n, err := html.Parse(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	return NewArena(*n), nil
}

func NewArena(root html.Node) *Arena {
	result := NewArenaRoot()
	result.transform(0, root)
	return result
}

func (a *Arena) Append(root html.Node) {
	a.transform(0, root)
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

	if currentId != node_index {
		a.AddChild(node_index, currentId)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.transform(currentId, *c)
	}
}

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

func (a *Arena) Stringify(nodeId int) string {
	n := a.Get(nodeId)
	res := n.String() + "\n"
	for _, c := range n.Children {
		res += "  " + a.Stringify(c)
	}
	return res
}

func (a *Arena) StringifyInformation(nodeId int) string {
	n := a.Get(nodeId)
	var res string
	if n.Type == html.TextNode {
		res = strings.TrimSpace(n.Data) + "\n"
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		res = n.GetAttr("src") + "\n"
	}

	for _, c := range n.Children {
		res += "  " + a.StringifyInformation(c)
	}
	return res
}

func (a *Arena) Rate(nodeId int) int {
	r := nodePoints
	n := a.Get(nodeId)
	for _, attr := range n.Attr {
		// doubled for compatibility
		//r += (attrKeyPoints + len(attr.Val)*attrValPoints) * 2
		r += attrKeyPoints + len(attr.Val)*attrValPoints
	}

	for _, c := range n.Children {
		r += a.Rate(c)
	}
	return r
}

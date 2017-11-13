// classify project classify.go
package classify

import (
	"strings"

	"golang.org/x/net/html"
)

// equality points
const textPoints int = 5
const nodePoints int = 5
const attrKeyPoints int = 1
const attrValPoints int = 1
const classPoints int = 2

type Node struct {
	Type html.NodeType
	Data string
	Attr []html.Attribute

	Children []int
	Parent   int
}

func NewNode(n html.Node) *Node {
	return &Node{
		Type: n.Type,
		Data: n.Data,
		Attr: n.Attr,
	}
}

func (n *Node) isInformative() bool {
	if strings.ToLower(n.Data) != "script" && strings.ToLower(n.Data) != "noscript" && strings.ToLower(n.Data) != "style" {
		if n.Type == html.TextNode {
			if len(strings.Trim(n.Data, "\x0d\x0a\x20\x09")) > 0 {
				return true
			}
		} else {
			if n.Type == html.ElementNode {
				if strings.ToLower(n.Data) == "img" {
					return true
				}
			}
		}
	}
	return false
}

/*
func (a *Arena) hasAnsector(node int, ansector int) bool {
	if node == ansector {
		return true
	} else {
		if a.list[node].Parent != 0 {
			return a.hasAnsector(a.list[node].Parent, ansector)
		}
	}
	return false
}
*/

// Function returns list of "informative" endings
func (a *Arena) GetInformative(nId int) []int {
	var r []int
	if a.List[nId].isInformative() {
		r = append(r, nId)
		return r
	}

	for _, id := range a.List[nId].Children {
		r = append(r, a.GetInformative(id)...)
	}

	return r
}

func (n *Node) Classes() []string {
	for _, a := range n.Attr {
		if a.Key == "class" {
			return strings.Fields(a.Val)
		}
	}

	return nil
}

func (n *Node) AddClass(className string) {
	for _, a := range n.Attr {
		if a.Key == "class" {
			if len(a.Val) > 0 {
				a.Val += " " + className
				return
			}
			a.Val = className
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: "class", Val: className})
}

func (n *Node) HasClass(c string) bool {
	for _, cl := range n.Classes() {
		if cl == c {
			return true
		}
	}
	return false
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

func CmpShallow(n1, n2 Node) *CmpResult {
	if n1.Type == n2.Type {
		if (n1.Type == html.TextNode) || (n1.Type == html.CommentNode) {
			r := CmpText(n1, n2)
			return &r
		}

		if n1.Data == n2.Data {
			r := &CmpResult{nodePoints, nodePoints}
			//return tagEqualityRate + (CmpAttr(n1.Attr, n2.Attr) * (1 - tagEqualityRate))
			attrResult := CmpAttr(n1.Attr, n2.Attr)
			if attrResult != nil {
				r.Append(*attrResult)
			}
			return r
		}
	}
	return nil
}

func CmpText(t1, t2 Node) CmpResult {
	if t1.Data == t2.Data {
		return CmpResult{nodePoints + textPoints, nodePoints + textPoints}
	}
	return CmpResult{nodePoints, nodePoints + textPoints}
}

func CmpAttr(attr1, attr2 []html.Attribute) *CmpResult {
	if len(attr1) == 0 && len(attr2) == 0 {
		return nil
	}

	r := &CmpResult{0, 0}

	if len(attr1) > 0 && len(attr2) > 0 {
		//var r, r1, r2 int
		for _, a1 := range attr1 {
			for _, a2 := range attr2 {
				if a1.Key == a2.Key && a1.Key != "class" && a2.Key != "class" {
					r.Sum += attrKeyPoints
					valueRate := CmpStrings(a1.Val, a2.Val)
					valueRate.Count *= attrValPoints
					valueRate.Sum *= attrValPoints
					r.Append(valueRate)
				}
			}
		}

		classes1 := []string{}
		classes2 := []string{}
		for _, a1 := range attr1 {
			if a1.Key == "class" {
				classes1 = strings.Fields(a1.Val)
				r.Count += len(classes1) * classPoints
			} else {
				r.Count += attrKeyPoints + attrValPoints*len(a1.Val)
			}
		}
		for _, a2 := range attr2 {
			if a2.Key == "class" {
				classes2 = strings.Fields(a2.Val)
				r.Count += len(classes2) * classPoints
			} else {
				r.Count += attrKeyPoints + attrValPoints*len(a2.Val)
			}
		}

		for _, c1 := range classes1 {
			for _, c2 := range classes2 {
				if c1 == c2 {
					r.Sum += classPoints
				}
			}
		}

		//double?
		r.Sum *= 2

		return r
	}

	if len(attr1) > 0 {
		r.Count = len(attr1)
	}

	if len(attr2) > 0 {
		r.Count = len(attr2)
	}

	return r
}

func (n Node) String() string {
	if n.Type == html.ElementNode {
		return "/" + n.Data + "[" + n.printAttr() + "]"
	}
	if n.Type == html.TextNode {
		return "text : " + n.Data
	}
	return n.Data
}

func (a *Arena) Path(n int) string {
	parent := a.Get(n).Parent
	if parent > 0 {
		return a.Path(parent) + a.Get(n).String()
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
		res = n.Data + "\n"
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		res = n.GetAttr("src") + "\n"
	}

	for _, c := range n.Children {
		res += "  " + a.StringifyInformation(c)
	}
	return res
}

func (n Node) printAttr() string {
	res := ""
	for _, a := range n.Attr {
		res += a.Key + "='" + a.Val + "', "
	}
	return res
}

func (n Node) GetAttr(key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

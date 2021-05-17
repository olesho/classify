// classify project classify.go
package arena

import (
	"fmt"
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
	Id       int

	// additional
	Ext interface{}
	//Volume float32
}

func NewNode(n html.Node, id int) *Node {
	return &Node{
		Type: n.Type,
		Data: n.Data,
		Attr: n.Attr,
		Id:   id,
	}
}

func (n *Node) Clone() *Node {
	c := &Node{
		Type:   n.Type,
		Data:   n.Data,
		Parent: n.Parent,
	}
	c.Attr = make([]html.Attribute, len(n.Attr))
	copy(c.Attr, n.Attr)
	c.Children = make([]int, len(n.Children))
	copy(c.Children, n.Children)
	return c
}

func (n *Node) isWholesome() bool {
	if (n.Type == html.TextNode) && (strings.TrimSpace(n.Data) != "") {
		return true
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				return true
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return true
			}
		}
	}
	return false
}

func (n *Node) Classes() []string {
	for _, a := range n.Attr {
		if a.Key == "class" {
			ss := strings.Fields(a.Val)
			result := []string{}
			for _, s := range ss {
				if !sliceContains(result, s) {
					result = append(result, s)
				}
			}
			return result
		}
	}

	return nil
}

func (n *Node) AddClass(className string) {
	for i, _ := range n.Attr {
		if n.Attr[i].Key == "class" {
			if len(n.Attr[i].Val) > 0 {
				n.Attr[i].Val += " " + className
				return
			}
			n.Attr[i].Val = className
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

// return textual representation for debugging
func (n Node) String() string {
	if n.Type == html.ElementNode {
		attrStr := ""
		if len(n.Attr) > 0 {
			attrStr = "[" + n.printAttr() + "]"
		}
		return "/" + n.Data + attrStr
	}
	if n.Type == html.TextNode {
		return "text : " + n.Data
	}
	return n.Data
}

func (n Node) printAttr() string {
	var attr []string
	for _, a := range n.Attr {
		attr = append(attr, fmt.Sprintf("%s='%s'", a.Key, a.Val))
	}
	return strings.Join(attr, ", ")
}

func (n Node) GetAttr(key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func sliceContains(sl []string, s string) bool {
	for _, item := range sl {
		if item == s {
			return true
		}
	}
	return false
}

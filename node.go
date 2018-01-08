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

	// used for multiple values
	DataArray []string
}

func NewNode(n html.Node) *Node {
	var dataArray []string
	if (n.Type == html.TextNode) && (strings.TrimSpace(n.Data) != "") {
		dataArray = []string{n.Data}
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				dataArray = []string{attr.Val}
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				dataArray = []string{attr.Val}
			}
		}
	}
	return &Node{
		Type:      n.Type,
		Data:      n.Data,
		Attr:      n.Attr,
		DataArray: dataArray,
	}
}

func (n *Node) isInformative() bool {
	if strings.ToLower(n.Data) != "script" && strings.ToLower(n.Data) != "noscript" && strings.ToLower(n.Data) != "style" {
		if n.Type == html.TextNode {
			//if len(strings.Trim(n.Data, "\x0d\x0a\x20\x09")) > 0 {
			if len(strings.TrimSpace(n.Data)) > 0 {
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

func (n *Node) Classes() []string {
	for _, a := range n.Attr {
		if a.Key == "class" {
			return strings.Fields(a.Val)
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

func (n Node) String() string {
	if n.Type == html.ElementNode {
		return "/" + n.Data + "[" + n.printAttr() + "]"
	}
	if n.Type == html.TextNode {
		return "text : " + n.Data
	}
	return n.Data
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

package sequence

import (
	"fmt"
	"strings"

	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
)

type Table struct {
	Arena         *arena.Arena
	TemplateArena *arena.Arena
	Members       []*arena.Node
	Rate          float32
	Volume        float32
	FieldSets     []FieldSet
}

type Field struct {
	Type int
	Val  string
}

type FieldSet struct {
	Type string
	Node *arena.Node

	Content []string
	IDs     []int
}

func (f FieldSet) String() string {
	s := fmt.Sprintf("type:%v ", f.Type)
	for i, f := range f.Content {
		s += fmt.Sprintf("%v:'%v' ", i, f)
	}
	return s
}

type Cell struct {
	Index int
	Rate  float32
}

func (c *Table) TemplateVolume() float32 {
	var vol float32 = .0
	for _, row := range c.FieldSets {
		switch row.Type {
		case "text":
			vol += textsVolume(row.Content)
		case "link":
			vol += linksVolume(row.Content)
		case "image":
			vol += imgsVolume(row.Content)
		}
	}
	return vol
}

func uniform(strs []string) bool {
	for _, s := range strs[1:] {
		if s != strs[0] {
			return false
		}
	}
	return true
}

func textsVolume(strs []string) float32 {
	var r float32 = .0
	for _, s := range strs {
		r += float32(len(s))
	}
	return r

	//smallest := float32(len(strs[0]))
	//for _, s := range strs[1:] {
	//	val := float32(len(s))
	//	if val < smallest {
	//		smallest = val
	//	}
	//}
	//return smallest * float32(len(strs))
}

func linksVolume(strs []string) float32 {
	var r float32 = .0
	for _, s := range strs {
		if len(s) > 0 {
			r += 0.1
		}
	}
	return r
}

func imgsVolume(strs []string) float32 {
	var r float32 = .0
	for _, s := range strs {
		if len(s) > 0 {
			r += 0.1
		}
	}
	return r
}

// WholesomeGroupFields checks each template arena node for valid field info; if found checks whole group of fields
func (c *Table) WholesomeGroupFields() []FieldSet {
	result := make([]FieldSet, 0)
	for _, n := range c.TemplateArena.List {
		if _, ok := WholesomeInfo(n); ok {
			ids := n.Ext.(*Additional).GroupIds
			if len(ids) == len(c.Members) {
				if values := extractFields(c.Arena, ids); values != nil {
					elementType := ""
					if n.Type == html.TextNode {
						elementType = "text"
					} else if n.Type == html.ElementNode {
						if n.Data == "a" {
							elementType = "link"
						} else if n.Data == "img" {
							elementType = "image"
						}
					}

					fieldSet := FieldSet{
						Type:    elementType,
						Content: values,
						Node:    n,
						IDs:     ids,
					}

					result = append(result, fieldSet)
				}
			}
		}
	}
	return result
}

// extractFields extracts all field values for certain type
func extractFields(arena *arena.Arena, ids []int) []string {
	content := make([]string, len(ids))
	for i, id := range ids {
		val, ok := WholesomeInfo(arena.Get(id))
		if !ok {
			return nil
		}
		content[i] = val
	}
	if !uniform(content) {
		return content
	}
	return nil
}

// WholesomeInfo extracts field value and type
func WholesomeInfo(n *arena.Node) (string, bool) {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data), true
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				return attr.Val, true
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return attr.Val, true
			}
		}
	}
	return "", false
}

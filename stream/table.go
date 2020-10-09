package stream

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
	Table         []Field
}

const MAX_TRIES = 4

const (
	NotField = iota
	TextField
	LinkField
	ImageField
)

type Field struct {
	Type    int
	Content []string
}

func (f Field) String() string {
	s := ""
	switch f.Type {
	case TextField:
		s += "type:text\n"
	case LinkField:
		s += "type:link\n"
	case ImageField:
		s += "type:image\n"
	}
	for i, f := range f.Content {
		s += fmt.Sprintf("%v: %v\n", i, f)
	}
	return s
}

type Cell struct {
	Index int
	Rate  float32
}

func (c *Table) TemplateVolume() float32 {
	var vol float32 = .0
	for _, row := range c.Table {
		switch row.Type {
		case TextField:
			vol += textsVolume(row.Content)
		case LinkField:
			vol += linksVolume(row.Content)
		case ImageField:
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

func (c *Table) WholesomeGroupTable() []Field {
	result := make([]Field, 0)
	for _, n := range c.TemplateArena.List {
		if _, fieldType := WholesomeInfo(n); fieldType != NotField {
			ids := n.Ext.(*Additional).GroupIds
			if len(ids) == len(c.Members) {
				if values := extractFields(c.Arena, ids, fieldType); values != nil {
					result = append(result, *values)
				}
			}
		}
	}
	return result
}

func extractFields(arena *arena.Arena, ids []int, fieldType int) *Field {
	values := &Field{}
	values.Content = make([]string, len(ids))
	for i, id := range ids {
		values.Content[i], values.Type = WholesomeInfo(arena.Get(id))
		if values.Type != fieldType {
			return nil
		}
	}
	if !uniform(values.Content) {
		return values
	}
	return nil
}

func WholesomeInfo(n *arena.Node) (string, int) {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data), TextField
	}
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				return attr.Val, ImageField
			}
		}
	}
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return attr.Val, LinkField
			}
		}
	}
	return "", NotField
}

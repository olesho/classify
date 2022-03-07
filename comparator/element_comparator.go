package comparator

import (
	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
)

type ElementComparator struct {
	a *arena.Arena
}

func NewElementComparator(a *arena.Arena) *ElementComparator {
	return &ElementComparator{ a }
}

func cmpAttr(n1, n2 *arena.Node) float32 {
	var total float32
	for _, attr1 := range n1.Attr {
		total += 1
		if attr1.Key == "class" {
			total += float32(len(n1.Classes()))
		} else {
			total += 1
		}
		//total += float32(len(attr1.Val))
	}
	for _, attr2 := range n2.Attr {
		total += 1
		if attr2.Key == "class" {
			total += float32(len(n2.Classes()))
		} else {
			total += 1
		}
		//total += float32(len(attr2.Val))
	}

	var coincided float32
	for _, attr1 := range n1.Attr {
		for _, attr2 := range n2.Attr {
			if attr1.Key == attr2.Key {
				coincided += 1
				if attr1.Key == "class" {
					classes1 := n1.Classes()
					classes2 := n2.Classes()
					for _, c1 := range classes1 {
						if hasStr(c1, classes2) {
							coincided += 1
						}
					}
				} else {
					coincided += cmpStrings(attr1.Val, attr2.Val)
				}
			}
		}
	}
	return coincided * 2 / total
}

func (c *ElementComparator) Cmp(idx1, idx2 int) float32 {
	n1, n2 := c.a.Get(idx1), c.a.Get(idx2)
	if n1.Type == n2.Type && n1.Type == html.TextNode {
		//return 1
		return cmpStrings(n1.Data, n2.Data)
	}
	if n1.Type == n2.Type {
		if n1.Data == n2.Data {
			return cmpAttr(n1, n2) + 1
		}
	}

	return 0
}

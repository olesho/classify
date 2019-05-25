// go-words project words.go
package html

import (
	"strings"

	"github.com/olesho/classify"
	"github.com/olesho/classify/go-words"
	"golang.org/x/net/html"
	"fmt"
)

type HtmlProcessor struct {
	*words.Processor
}

func NewHtmlProcessor() *HtmlProcessor {
	return &HtmlProcessor{words.NewProcessor()}
}

type Node classify.Node
type Word words.Word

func (c1 Node) SameKind(c2 words.Comparable) bool {
	cmp := c2.(Node)
	if c1.Type != cmp.Type {
		return false
	}

	if c1.Type == html.ElementNode {
		if c1.Data != cmp.Data {
			return false
		}
	}

	return true
}

func (w Word) String() string {
	nn := w.toNodes()
	res := make([]string, len(w))
	for i, n := range nn {
		if n.Type == html.ElementNode {
			res[i] = n.Data
		} else if n.Type == html.TextNode {
			res[i] = "TEXT"
		} else {
			res[i] = "OTHER"
		}
	}
	return strings.Join(res, " ")
}

func (w Word) toNodes() []Node {
	res := make([]Node, len(w))
	for i, _ := range w {
		res[i] = w[i].(Node)
	}
	return res
}

func FindValues(word, text Word) [][]string {
	positions := words.FindPositions(words.Word(word), words.Word(text))
	fieldsTable := make([][]string, len(positions))
	for i, p := range positions {
		fmt.Println(text[p])
		fieldsTable[i] = extractValues(text[p: p+len(word)])
	}
	return fieldsTable
}

func extractValues(w Word) []string {
	fields := []string{}
	for _, c := range w {
		if char, ok := c.(Node); ok {
			if field, isInfo := classify.Node(char).Info(); isInfo {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

func FilterValues(fieldsTable [][]string) [][]string {
	result := make([][]string, len(fieldsTable))
	if len(fieldsTable) > 0 {
		for i, _ := range fieldsTable[0] {
			if !isFieldMonotone(i, fieldsTable) {
				for j, _ := range result {
					result[j] = append(result[j], fieldsTable[j][i])
				}
			}
		}
	}
	return result
}

func isFieldMonotone(fieldIndex int, fieldsTable [][]string) bool {
	val := fieldsTable[0][fieldIndex]
	for _, fields := range fieldsTable {
		if fields[fieldIndex] != val {
			return false
		}
	}
	return true
}
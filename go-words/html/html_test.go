// go-words project words.go
package html

import (
	"log"
	"net/http"
	"testing"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
	"fmt"
)

func TestArenaVocabulary(t *testing.T) {
	//	a := assert.New(t)
	//resp, err := http.Get("http://www.bbc.com/")
	resp, err := http.Get("http://news.ycombinator.com/")
	if err != nil {
		log.Println(err)
		return
	}
	if resp != nil {
		defer resp.Body.Close()
		n, err := html.Parse(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		arena := classify.NewArena(*n)
		cl := NewHtmlProcessor()
		for _, n := range arena.List {
			cl.Next(Node(n))
		}

		cl.SortVocabulary()
		for i, _ := range cl.Vocabulary() {
			fmt.Println(Word(cl.Vocabulary()[i].Word))
			fmt.Println(cl.Count(cl.Vocabulary()[i].Word))
		}
	}
}

func Tags(nn []classify.Node) []string {
	res := make([]string, len(nn))
	for i, n := range nn {
		if n.Type == html.ElementNode {
			res[i] = n.Data
		} else if n.Type == html.TextNode {
			res[i] = "TEXT"
		} else {
			res[i] = "OTHER"
		}
	}
	return res
}

func (w Word) Tags() []string {
	nn := wordToNodes(w)
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
	return res
}

func wordToNodes(w Word) []Node {
	res := make([]Node, len(w))
	for i, _ := range w {
		res[i] = w[i].(Node)
	}
	return res
}
package html

import (
	"io"
	"golang.org/x/net/html"
	"log"
	"github.com/olesho/classify"
)

type Engine struct {
	processor *HtmlProcessor
}

func NewEngine() *Engine {
	return &Engine{
		processor: NewHtmlProcessor(),
	}
}

func (e *Engine) Train(r io.Reader) {
	n, err := html.Parse(r)
	if err != nil {
		log.Println(err)
		return
	}

	arena := classify.NewArena(*n)
	for _, n := range arena.List {
		e.processor.Next(Node(n))
	}
}

func (e *Engine) Sort() {
	e.processor.SortVocabulary()
}

func (e *Engine) Parse(r io.Reader, index int) [][]string {
	bestWord := Word(e.processor.Vocabulary()[index].Word)
	text := Word(e.processor.Text)
	return FilterValues(FindValues(bestWord, text))
}
// classify project arena_test.go
package deep

import (
	"os"
	"testing"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

func TestDeepClassify(t *testing.T) {
	f, err := os.Open("../examples/BBC - Homepage.html")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	n, err := html.Parse(f)
	if err != nil {
		t.Error(err)
	}
	a := classify.NewArena(*n)
	c := NewDeepClassificator(a)
	c.Run()

	/*
		for _, b := range c.bags.List {
			for _, i := range b.Content {
				t.Log(c.Get(i).String())
			}
		}
	*/
	if len(c.bags.List) < 1 {
		t.Error("Cant be 0 result")
	}
}

func TestPattern(t *testing.T) {
	f, err := os.Open("../examples/BBC - Homepage.html")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	n, err := html.Parse(f)
	if err != nil {
		t.Error(err)
	}
	a := classify.NewArena(*n)
	c := NewDeepClassificator(a)
	c.Run()
	pattern := classify.Pattern(&c.Arena, c.Bags().List[0])
	t.Log(pattern.Stringify(1))
}

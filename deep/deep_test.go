// classify project arena_test.go
package deep

import (
	"os"
	"testing"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

func TestBBC(t *testing.T) {
	//f, err := os.Open("../examples/I went full nomad and it (almost) broke me – Hacker Noon.html")
	//f, err := os.Open("../examples/BBC - Homepage.html")
	f, err := os.Open("../examples/Hacker News.html")
	//f, err := os.Open("../examples/Hacker Noon.html")
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

	for bIndex, bag := range c.bags.List {
		t.Log("Bag size:", len(bag.Content), "Bag rate:", bag.Rate, "Bag sum:", bag.Sum, "Bag Index:", bag.Index)
		//t.Log(bag.Content)
		for _, i := range bag.Content {
			t.Log("Item:", i)
			t.Log(c.StringifyInformation(i))
		}
		t.Log("=================================================================================================")
		if bIndex == 10 {
			break
		}
	}

	/*
		bag := c.bags.List[0]
		for _, i := range bag.Content {
			t.Log(c.Path(i))
		}

		path := classify.GeneratePath(&c.Arena, bag)
		t.Log(path)
		pattern := classify.GeneratePattern(&c.Arena, bag)
		t.Log(pattern.PrintList())
	*/
}

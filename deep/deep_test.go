// classify project arena_test.go
package deep

import (
	"os"
	"testing"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

func TestBBC(t *testing.T) {
	//f, err := os.Open("../examples/BBC - Homepage.html")
	//f, err := os.Open("../examples/Hacker News.html")
	f, err := os.Open("../examples/Hacker Noon.html")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	n, err := html.Parse(f)
	if err != nil {
		t.Error(err)
	}
	a := classify.NewArena(*n)

	//ids := a.FindNodeIdByAttr("class", "media-list__item media-list__item--4")
	//t.Log(ids)

	c := NewDeepClassificator(a)
	c.Run()

	//bag := c.bags.List[3]
	//t.Log("Bag size:", len(bag.Content), "Bag rate:", bag.Rate)

	for bIndex, bag := range c.bags.List {
		t.Log("Bag:", bIndex, "Bag size:", len(bag.Content), "Bag rate:", bag.Rate, "Bag sum:", bag.Sum)
		t.Log(bag.Content)
		for _, i := range bag.Content {
			//t.Log(c.Get(i).String())
			t.Log(c.StringifyInformation(i))
		}
		t.Log("=================================================================================================")
		if bIndex == 10 {
			break
		}
	}

	//path := classify.GeneratePath(&c.Arena, bag)
	//t.Log(path)
	//pattern := classify.GeneratePattern(&c.Arena, bag)
	//t.Log(pattern.PrintList())
}

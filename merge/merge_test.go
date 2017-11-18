// classify project merge_test.go
package merge

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

	c := NewMergeClassificator(a)
	c.Run()

	for bIndex, bag := range c.bags.List {
		t.Log("Bag:", bIndex, "Bag size:", len(bag.Content), "Bag rate:", bag.Efficacy())
		//t.Log(bag.Content)
		for _, i := range bag.Content {
			//t.Log(c.Get(i).String())
			t.Log("Item:", i)
			t.Log(c.StringifyInformation(i))
		}
		t.Log("=================================================================================================")
		if bIndex == 10 {
			break
		}
	}
}

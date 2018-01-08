// classify project merge_test.go
package merge

import (
	"os"
	"testing"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

/*
func TestHackerNoon(t *testing.T) {
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

	nodes := c.Arena.IndexesByClass("collectionHeader-navItem")
	bags := c.BagsContaining(nodes)
	if len(bags) < 1 {
		t.Error("Predicted bag #1 not found!")
		return
	}

	if len(bags[0].Content) != 4 {
		t.Error("Bag #1 size not equal to predicted")
		return
	}

	for bIndex, bag := range c.bags.List {
		t.Log("Bag:", bIndex, "Bag size:", len(bag.Content), "Bag rate:", bag.Efficacy())
		//t.Log(bag.Content)
		for _, i := range bag.Content {
			//t.Log(c.Get(i).String())
			t.Log("Item:", i)
			//t.Log(c.Get(i).String())
		}
		t.Log("=================================================================================================")
		if bIndex == 10 {
			break
		}
	}
}
*/

func TestHackerNews(t *testing.T) {
	f, err := os.Open("../examples/Hacker News.html")
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

	for bIndex, bag := range c.Bags().List {
		t.Log("Bag:", bIndex, "Bag size:", len(bag.Content), "Bag rate:", bag.Efficacy())
		//t.Log(bag.Content)
		for _, i := range bag.Content {
			//t.Log("Item:", i)
			t.Log(c.StringifyInformation(i))
		}
		t.Log("=================================================================================================")
		if bIndex == 10 {
			break
		}
	}

}

/*
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

	t.Log(c.Arena.Stringify(65))

	nodes := c.Arena.IndexesByClass("collectionHeader-navItem")
	bags := c.BagsContaining(nodes)
	if len(bags) < 1 {
		t.Error("Predicted bag #1 not found!")
		return
	}
	if len(bags[0].Content) != 3 {
		t.Error("Bag #1 size not equal to predicted")
		return
	}

	nodes = c.Arena.IndexesByAttr("class", "u-size8of12 u-xs-size12of12 u-minHeight400 u-xs-height350 u-overflowHidden js-trackedPost u-relative u-imageSpectrum")
	for i, _ := range nodes {
		nodes[i] = c.Arena.Get(nodes[i]).Parent
	}
	t.Log(nodes)
	bags = c.BagsContaining(nodes)
	if len(bags) < 1 {
		t.Error("Predicted bag #2 not found!")
		return
	}
	t.Log(bags)
	if len(bags[0].Content) != 3 {
		t.Error("Bag #2 size not equal to predicted")
	}

	/*
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
*/
//}

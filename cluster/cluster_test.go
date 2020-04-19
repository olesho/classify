package cluster

import (
	"bufio"
	"fmt"
	"github.com/olesho/classify"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"os"
	"testing"
)

//func TestRenderLabels(t *testing.T) {
//	a := assert.New(t)
//	f, _ := os.Open("examples/bbc.html")
//	defer f.Close()
//	reader := bufio.NewReader(f)
//	n, err := html.Parse(reader)
//	a.NoError(err)
//
//	arena := classify.NewArena(*n)
//	for i, el := range arena.List {
//		arena.List[i].Attr = append(el.Attr, html.Attribute{
//			Key: "arid",
//			Val: fmt.Sprint(el.Id),
//		})
//	}
//	text, _ := arena.RenderString(0)
//
//	fmt.Println(text)
//}

func TestYcomb(t *testing.T) {
	a := assert.New(t)

	//f, _ := os.Open("examples/fb.html")
	f, _ := os.Open("examples/ycomb.html")
	//f, _ := os.Open("examples/hackernoon.html")
	//f, _ := os.Open("examples/pravda.html")
	//f, _ := os.Open("examples/bbc.html")
	//f, _ := os.Open("examples/cnn.html")
	//f, _ := os.Open("examples/test1.html")
	//f, _ := os.Open("examples/test2.html")
	defer f.Close()
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	a.NoError(err)

	arena := classify.NewArena(*n)
	series := Extract(arena).Matrix[0]

	for _, c := range series.Group.Clusters {
		v := c.TemplateVolume()
		if v > 0 {
			fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::::::::::")
			fmt.Println(c.TemplateArena.RenderString(0))
			fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::::::::::")
			for _, ss := range c.Table {
				fmt.Println(ss)
			}
			fmt.Printf("similarity rate: %v, template volume: %v\n", c.Rate, v)
			fmt.Println("___________________________________________________")
		}
	}

	template := series.Nonuniform().Patterns()
	for _, r := range template.Chains {
		fmt.Println(r.XPath())
	}
	fmt.Printf("size: %v, volume: %v, group volume: %v\n", series.Group.Size, series.Group.Volume, series.Group.GroupVolume)
	return
}

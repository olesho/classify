package cluster

import (
	"fmt"
	"strings"
	"testing"

	"github.com/olesho/classify/arena"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

//func TestRenderLabels(t *testing.T) {
//	a := assert.New(t)
//	f, _ := os.Open("examples/bbc.html")
//	defer f.Close()
//	reader := bufio.NewReader(f)
//	n, err := html.Parse(reader)
//	a.NoError(err)
//
//	arena := arena.NewArena(*n)
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

	n1, err := html.Parse(strings.NewReader(`
<html>
	<body>
		<div>
			<p>Hello 1</p>
		</div>
		<div>
			<p>Hello 2</p>
		</div>
		<div>
			<p>Hello 3</p>
		</div>
	</body>
</html>
`))
	a.NoError(err)

	n2, err := html.Parse(strings.NewReader(`
<html>
	<body>
		<div>
			<p>Hello 4</p>
		</div>
		<div>
			<p>Hello 5</p>
		</div>
		<div>
			<p>Hello 6</p>
		</div>
	</body>
</html>
`))
	a.NoError(err)

	arena := arena.NewArena()
	arena.Append(*n1)
	arena.Append(*n2)

	//s, _ := arena.RenderString(0)
	//fmt.Println(s)

	for i, item := range arena.List {
		if item.Data == "div" {
			fmt.Println(i)
		}
	}

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
	fmt.Printf("total chains: %v\n", len(template.Chains))
	for _, r := range template.Chains {
		fmt.Println(r.XPath())
	}
	fmt.Printf("size: %v, volume: %v, group volume: %v\n", series.Group.Size, series.Group.Volume, series.Group.GroupVolume)
	return
}

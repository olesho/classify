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

func TestYcomb(t *testing.T) {
	a := assert.New(t)

	//f, _ := os.Open("examples/ycomb.html")
	f, _ := os.Open("examples/hackernoon.html")
	//f, _ := os.Open("examples/pravda.html")
	//f, _ := os.Open("examples/bbc.html")
	defer f.Close()
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	a.NoError(err)

	arena := classify.NewArena(*n)

	//for i, n := range arena.List {
	//	if n.HasClass("rank") {
	//		fmt.Println(i)
	//	}
	//	//fmt.Println(strings.Replace(arena.Chain(n.Id, 0).XPath(), "\n", " ", -1))
	//}

	series, template := Extract2(arena).BestPattern()

	//template := series.Nonuniform().Informative().Patterns()
	//template := series.Patterns()
	for _, r := range template.Chains {
		fmt.Println(r.XPath())
	}
	fmt.Println(len(series.Matrix))
	return

	//for _, row := range series.Nonuniform().Matrix {
	//	for _, n := range row {
	//		//str, _ := arena.RenderString(n.Id)
	//		str := arena.StringifyInformation(n.Id)
	//		fmt.Println(str)
	//	}
	//	fmt.Println("__________________________________________________________________________________________________")
	//}
}

package bags

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/olesho/classify"
	"fmt"
	"bufio"
	"golang.org/x/net/html"
)

func TestProcessor(t *testing.T) {
	a := assert.New(t)

	f, _ := os.Open("test.html")
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	a.NoError(err)

	arena := classify.NewArena(*n)

	nodes, err := Parse(arena, 0)
	a.NoError(err)

	for _, n := range nodes {
		//fmt.Println(arena.RenderString(n.Id))
		fmt.Println(arena.StringifyInformation(n.Id))
		fmt.Println()
	}
	fmt.Println("===============================================")


}

//func TestComparator(t *testing.T) {
//	a := assert.New(t)
//
//	data := `
//	<html>
//		<head></head>
//		<body>
//			<div id="block_1">
//				<h1>Test header 1</h1>
//				<p class="super clearfix test">Some text</p>
//			</div>
//			<div id="block_2">
//				<h1>Test header 2</h1>
//				<p class="super clearfix test">Another text</p>
//			</div>
//			<img src="/i.png"></img>
//		</body>
//	</html>
//	`
//
//	ar1, err := classify.NewArenaHtml(data)
//	a.NoError(err)
//
//	for _, n := range ar1.List {
//		fmt.Println(n.Id, n.String())
//	}
//
//	fmt.Println()
//	fmt.Println("StrictComparator:", StrictComparator(ar1, ar1.Get(7), ar1.Get(12)))
//	fmt.Println("elementComparator:", elementComparator(ar1, ar1.Get(7), ar1.Get(12)))
//	fmt.Println("cmpAttr:", cmpAttr(ar1.Get(7).Attr, ar1.Get(12).Attr))
//	fmt.Println("ColumnComparator:", ColumnComparator(ar1, ar1.Get(7), ar1.Get(12)))
//	fmt.Println("ChildComparator", ChildComparator(ar1, ar1.Get(7), ar1.Get(12)))
//	fmt.Println("ExtendedComparator", ExtendedComparator(ar1, ar1.Get(7), ar1.Get(12)))
//
//	fmt.Println()
//	fmt.Println("StrictComparator:", StrictComparator(ar1, ar1.Get(4), ar1.Get(9)))
//	fmt.Println("elementComparator:", elementComparator(ar1, ar1.Get(4), ar1.Get(9)))
//	fmt.Println("cmpAttr:", cmpAttr(ar1.Get(4).Attr, ar1.Get(9).Attr))
//	fmt.Println("ColumnComparator:", ColumnComparator(ar1, ar1.Get(4), ar1.Get(9)))
//	fmt.Println("ChildComparator", ChildComparator(ar1, ar1.Get(4), ar1.Get(9)))
//	fmt.Println("ExtendedComparator", ExtendedComparator(ar1, ar1.Get(4), ar1.Get(9)))
//}
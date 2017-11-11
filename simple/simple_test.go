// classify project simple_test.go
package simple

import (
	"strings"
	"testing"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

func TestClassify(t *testing.T) {
	n, _ := html.Parse(strings.NewReader(`
	<html>
		<head></head>
		<body>
			<div id="block_1">
				<h1>Test header 1</h1>
				<p class="super clearfix test">Some text</p>
			</div>
			<div id="block_2">
				<h1>Test header 2</h1>
				<p class="super clearfix test">Another text</p>
			</div>
			<div id="block_3">
				<h1>Test header 3</h1>
			</div>
			<img src="/i.png"></img>
		</body>
	</html>
	`))
	a := classify.NewArena(*n)
	rate := a.CmpColumn(9, 12)
	t.Log(rate)
}

func TestCmpShallow(t *testing.T) {
	n, _ := html.Parse(strings.NewReader(`
	<html>
		<head></head>
		<body>
			<div id="block_1" class="box">
				<h1>Test header 1</h1>
				<p class="super clearfix test">Some text</p>
			</div>
			<div id="block_2" class="box">
				<h1>Test header 2</h1>
				<p class="super clearfix test">Another text</p>
			</div>
			<div id="block_3">
				<h1>Test header 3</h1>
			</div>
			<img id="img_1" src="/i.png"></img>
		</body>
	</html>
	`))
	a := classify.NewArena(*n)
	n1 := a.FindByAttr("id", "block_1")
	n2 := a.FindByAttr("id", "block_3")
	rate := classify.CmpShallow(*n1, *n2)
	t.Log(rate)
}

func TestCmpColumn(t *testing.T) {
	n, _ := html.Parse(strings.NewReader(`
	<html>
		<head></head>
		<body>
			<div id="block_1" class="box">
				<h1>Test header 1</h1>
				<p class="super clearfix test">Some text</p>
				<i>
					<img src="p1.jpg" />
				</i>
			</div>
			<div id="block_2" class="box">
				<h1>Test header 2</h1>
				<p class="super clearfix test">Another text</p>
				<i>
					<img src="p2.jpg" />
				</i>
			</div>
			<div id="block_3">
				<h1>Test header 3</h1>
				<i>
					<img src="p3.jpg" />
				</i>
			</div>
			<img id="img_1" src="/i.png"></img>
		</body>
	</html>
	`))
	a := classify.NewArena(*n)
	n1 := a.FindNodeIdByAttr("src", "p1.jpg")
	n2 := a.FindNodeIdByAttr("src", "p2.jpg")
	rate := a.CmpColumn(n1, n2)
	t.Log(rate)
}

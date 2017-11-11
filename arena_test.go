// classify project arena_test.go
package classify

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestNewArena(t *testing.T) {
	n, _ := html.Parse(strings.NewReader(`
	<html>
		<head></head>
		<body>
			<div>
				<h1>Test header</h1>
				<p>Paragraph</p>
			</div>
		</body>
	</html>
	`))
	a := NewArena(*n)
	root := a.Get(a.Get(0).Children[1])
	if root.Data != "html" {
		t.Error("Error: no root 'html'")
	}

	child := a.Get(root.Children[1])
	if child.Data != "body" {
		t.Error("Error: no root 'body'")
	}

}

func TestGetInformative(t *testing.T) {
	n, _ := html.Parse(strings.NewReader(`
	<html>
		<head></head>
		<body>
			<div>
				<h1>Test header</h1>
				<p class="super clearfix test">Paragraph</p>
			</div>
			<img src="/i.png"></img>
		</body>
	</html>
	`))
	a := NewArena(*n)

	ids := a.GetInformative(1)

	if !a.List[a.List[ids[1]].Parent].HasClass("clearfix") {
		t.Error("HasClass not working")
	}
}

/*
func TestCmpLightRate(t *testing.T) {
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
			<img src="/i.png"></img>
		</body>
	</html>
	`))
	a := NewArena(*n)
	rate := a.CmpLightRate(*a.FindByAttr("id", "block_1"), *a.FindByAttr("id", "block_1"))
	t.Log(rate)
}

func TestCmpDeepRate(t *testing.T) {
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
			<img src="/i.png"></img>
		</body>
	</html>
	`))
	a := NewArena(*n)
	rate := a.CmpDeepRate(*a.FindByAttr("id", "block_1"), *a.FindByAttr("id", "block_1"))
	t.Log(rate)
}
*/

func TestCloneArena(t *testing.T) {
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
			<img src="/i.png"></img>
		</body>
	</html>
	`))
	a := NewArena(*n)
	b := NewArenaRoot()
	id := a.FindNodeIdByAttr("id", "block_1")
	a.Clone(id)
	t.Log(b.PrintList())
}

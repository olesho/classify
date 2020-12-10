// classify project arena_test.go
package arena

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	root := a.Get(a.Get(0).Children[0])
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

	ids := a.Wholesome(1)

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

func TestRenderArena(t *testing.T) {
	a := assert.New(t)

	data := `
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
	`

	ar1, err := NewArenaHtml(data)
	a.NoError(err)

	s, err := ar1.RenderString(0)
	a.NoError(err)
	a.Equal(s, `<html><head></head><body><div id="block_1"><h1>Test header 1</h1><p class="super clearfix test">Some text</p></div><div id="block_2"><h1>Test header 2</h1><p class="super clearfix test">Another text</p></div><img src="/i.png"/></body></html>`)
}

func TestCloneArena(t *testing.T) {
	a := assert.New(t)

	ar1, err := NewArenaHtml(`
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
	`)
	a.NoError(err)

	idx := ar1.IndexesByAttr("id", "block_1")[0]

	ar2 := ar1.CloneBranch(idx)
	s1, err := ar1.RenderString(idx)
	a.NoError(err)
	s2, err := ar2.RenderString(0)
	a.NoError(err)

	a.Equal(s1, s2)
}

func TestArenaStructure(t *testing.T) {
	a := assert.New(t)

	ar1, err := NewArenaHtml(`
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
	`)
	a.NoError(err)

	list := []string{"", "html", "head", "body", "div", "h1", "Test header 1", "p", "Some text", "div", "h1", "Test header 2", "p", "Another text", "img"}

	for i, el := range ar1.List {
		if el.Data != list[i] {
			a.Fail("Wrong arena structure")
		}
	}

}

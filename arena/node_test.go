// classify project classify.go
package arena

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestNode(t *testing.T) {
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
				<p id="specific" class="super clearfix test">Another text</p>
			</div>
			<div id="block_3">
				<h1>Test header 3</h1>
			</div>
			<img src="/i.png"></img>
		</body>
	</html>
	`))
	a := NewArena(*n)
	node := a.IndexesByAttr("id", "specific")
	arr := a.PathArray(node[0])
	if a.Get(arr[0]).Data != "p" {
		t.Error("Wrong item in path array")
	}
	if a.Get(arr[1]).Data != "div" {
		t.Error("Wrong item in path array")
	}
	if a.Get(arr[2]).Data != "body" {
		t.Error("Wrong item in path array")
	}
	if a.Get(arr[3]).Data != "html" {
		t.Error("Wrong item in path array")
	}
}

func TestAddClass(t *testing.T) {
	r := Node{
		Type: html.ElementNode,
		Attr: make([]html.Attribute, 0),
	}

	r.AddClass("one")
	r.AddClass("two")
	r.AddClass("three")

	if !r.HasClass("one") {
		t.Error("Class not added!")
	}
	if !r.HasClass("two") {
		t.Error("Class not added!")
	}
	if !r.HasClass("three") {
		t.Error("Class not added!")
	}
}

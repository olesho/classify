// classify project classify.go
package arena

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	a := assert.New(t)

	ar1, err := NewArenaHtml(`
	<html>
		<head></head>
		<body>
			<div id="block_1">
				<h1>Test header 1</h1>
				<p class="super clearfix test">Some text</p>
			</div>
			<div>
				<div id="block_2">
					<h1>Test header 2</h1>
					<p class="super clearfix test">Another text</p>
				</div>
			</div>
			<img src="/i.png"></img>
		</body>
	</html>
	`)
	a.NoError(err)

	idx := ar1.IndexesByAttr("class", "super clearfix test")
	c1 := ar1.Chain(ar1.Get(idx[0]).Children[0], 0)
	c2 := ar1.Chain(ar1.Get(idx[1]).Children[0], 0)

	t1 := &Template{[]Chain{c1}}
	t2 := &Template{[]Chain{c2}}

	mt := MergeTemplates(t1, t2)

	t.Log(mt)
}

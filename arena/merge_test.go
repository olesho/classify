// classify project arena_test.go
package arena

// func TestMergeNodes(t *testing.T) {
// 	n, _ := html.Parse(strings.NewReader(`
// 	<html>
// 		<head></head>
// 		<body>
// 			<div id="block_1">
// 				<h1>Test header 1</h1>
// 				<p class="super clearfix test">Some text</p>
// 			</div>
// 			<div id="block_2">
// 				<h1>Test header 2</h1>
// 				<p class="super clearfix test">Another text</p>
// 			</div>
// 			<div id="block_3">
// 				<h1>Test header 3</h1>
// 			</div>
// 			<img src="/i.png"></img>
// 		</body>
// 	</html>
// 	`))
// 	a := NewArena(*n)
// 	id1 := a.IndexesByAttr("id", "block_1")
// 	id2 := a.IndexesByAttr("id", "block_3")

// 	a1 := a.CloneBranch(id1[0])
// 	a2 := a.CloneBranch(id2[0])

// 	res := Merge(a1, a2, 0, 0)
// 	t.Log(res.PrintList())
// }

// func TestMergeInside(t *testing.T) {
// 	n, _ := html.Parse(strings.NewReader(`
// 	<html>
// 		<head></head>
// 		<body>
// 			<div id="block_1">
// 				<h1>Test header 1</h1>
// 				<p class="super clearfix test">Some text</p>
// 			</div>
// 			<div id="block_2">
// 				<h1>Test header 2</h1>
// 				<p class="super clearfix test">Another text</p>
// 			</div>
// 			<div id="block_3">
// 				<h1>Test header 3</h1>
// 			</div>
// 			<img src="/i.png"></img>
// 		</body>
// 	</html>
// 	`))
// 	a := NewArena(*n)
// 	id1 := a.IndexesByAttr("id", "block_1")
// 	id2 := a.IndexesByAttr("id", "block_3")

// 	res := Merge(a, a, id1[0], id2[0])
// 	t.Log(res.PrintList())
// }

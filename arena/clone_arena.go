// classify project classify.go
package arena

// func (a *Arena) Clone(srcId int) (dest *Arena) {
// 	dest = NewArenaRoot()
// 	lastId := a.lastNode(srcId) + 1

// 	dest.List = make([]*Node, lastId-srcId)
// 	src := a.List[srcId:lastId]
// 	for i, _ := range src {
// 		c := src[i].Clone()
// 		for j, _ := range c.Children {
// 			c.Children[j] = c.Children[j] - srcId
// 		}
// 		c.Parent = c.Parent - srcId
// 		dest.List[i] = c

// 	}
// 	return dest

// }

// // gets last node among descendants
// func (a *Arena) lastNode(id int) int {
// 	children := a.Get(id).Children
// 	size := len(children)
// 	if size == 0 {
// 		return id
// 	}
// 	return a.lastNode(children[size-1])
// }

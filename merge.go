// classify project classify.go
package classify

import (
	//	"fmt"

	"golang.org/x/net/html"
)

type Coord struct {
	i int
	j int
}

func Merge(a1, a2 *Arena, id1, id2 int) *Arena {
	m := make(map[int]map[int][]Coord)
	makePairs(a1, a2, id1, id2, m)
	resArena := NewArenaRoot()
	buildRoot(a1, a2, id1, id2, m, resArena)
	return resArena
}

func buildRoot(a1, a2 *Arena, id1, id2 int, pairs map[int]map[int][]Coord, dest *Arena) {
	resNode := MergeShallow(a1.Get(id1), a2.Get(id2))

	if resNode != nil {
		dest.List = append(dest.List, *resNode)

		/* save informatives
		if a2.Get(id2).isInformative() {

		}
		*/

		resId := len(dest.List) - 1
		for _, coord := range pairs[id1][id2] {
			build(a1, a2, a1.Get(id1).Children[coord.i], a2.Get(id2).Children[coord.j], pairs, dest, resId)
		}
	}
}

func build(a1, a2 *Arena, id1, id2 int, pairs map[int]map[int][]Coord, dest *Arena, destId int) {
	resNode := MergeShallow(a1.Get(id1), a2.Get(id2))
	if resNode != nil {
		dest.List = append(dest.List, *resNode)
		resId := len(dest.List) - 1
		dest.AddChild(destId, resId)

		/* save informatives
		if a2.Get(id2).isInformative() {

		}
		*/

		for _, coord := range pairs[id1][id2] {
			build(a1, a2, a1.Get(id1).Children[coord.i], a2.Get(id2).Children[coord.j], pairs, dest, resId)
		}
	}
}

func makePairs(a1, a2 *Arena, id1, id2 int, pairs map[int]map[int][]Coord) (r *CmpResult) {
	n1 := a1.List[id1]
	n2 := a2.List[id2]

	if n1.Type == n2.Type {
		if n1.Data == n2.Data {
			r = &CmpResult{nodePoints, nodePoints}
			//twin := NewFamily(&HtmlNode{Type: n1.Type, Data: n1.Data, Attr: make([]html.Attribute, 0)})

			// compare attributes (except classes)
			for _, a1 := range n1.Attr {
				for _, a2 := range n2.Attr {
					if a1.Key == a2.Key && a1.Key != "class" && a2.Key != "class" {
						r.Sum += attrKeyPoints
						r.Count += attrKeyPoints
						attr := html.Attribute{Key: a1.Key}
						if a1.Val == a2.Val {
							attr.Val = a1.Val
							r.Sum += len(a1.Val)
							r.Count += len(a2.Val)
						}
						//						twin.Attr = append(twin.Attr, attr)
					}
				}
			}

			//compare classes
			classes1 := n1.Classes()
			classes2 := n2.Classes()
			for _, c1 := range classes1 {
				for _, c2 := range classes2 {
					if c1 == c2 {
						r.Sum += classPoints
					}
				}
			}

			l1 := len(n1.Children)
			l2 := len(n2.Children)

			if (l1 > 0) && (l2 > 0) {
				//children := make([]int, 0)
				if _, ok := pairs[id1]; !ok {
					pairs[id1] = make(map[int][]Coord)
				}

				matrix := make([][]*CmpResult, l1)
				for i, c1 := range n1.Children {
					matrix[i] = make([]*CmpResult, l2)
					for j, c2 := range n2.Children {
						matrix[i][j] = makePairs(a1, a2, c1, c2, pairs)
					}
				}

				first_used := make([]bool, l1)
				second_used := make([]bool, l2)

				max_i, max_j, max_res := bestOfSquareMatrix(matrix, first_used, second_used)
				for max_res != nil {
					//					max_rate := max_res.Result()
					first_used[max_i] = true
					second_used[max_j] = true
					pairs[id1][id2] = append(pairs[id1][id2], Coord{max_i, max_j})
					r.Append(*max_res)

					max_i, max_j, max_res = bestOfSquareMatrix(matrix, first_used, second_used)
				}
			}
			return r
		}
	}
	return nil
}

/*
func MakePairsMock(a1, a2 *Arena, id1, id2 int, pairs map[int]map[int][]Coord) (r *CmpResult) {
	fmt.Println("Mock")
	n1 := a1.List[id1]
	n2 := a2.List[id2]

	if n1.Type == n2.Type {
		if n1.Data == n2.Data {
			r = &CmpResult{nodePoints, nodePoints}
			//twin := NewFamily(&HtmlNode{Type: n1.Type, Data: n1.Data, Attr: make([]html.Attribute, 0)})

			// compare attributes (except classes)
			for _, a1 := range n1.Attr {
				for _, a2 := range n2.Attr {
					if a1.Key == a2.Key && a1.Key != "class" && a2.Key != "class" {
						r.Sum += attrKeyPoints
						r.Count += attrKeyPoints
						attr := html.Attribute{Key: a1.Key}
						if a1.Val == a2.Val {
							attr.Val = a1.Val
							r.Sum += len(a1.Val)
							r.Count += len(a2.Val)
						}
						//						twin.Attr = append(twin.Attr, attr)
					}
				}
			}

			//compare classes
			classes1 := n1.Classes()
			classes2 := n2.Classes()
			for _, c1 := range classes1 {
				for _, c2 := range classes2 {
					if c1 == c2 {
						r.Sum += classPoints
					}
				}
			}

			l1 := len(n1.Children)
			l2 := len(n2.Children)

			if (l1 > 0) && (l2 > 0) {

				fmt.Println("Children:", l1, l2)

				//children := make([]int, 0)
				if _, ok := pairs[id1]; !ok {
					pairs[id1] = make(map[int][]Coord)
				}

				matrix := make([][]*CmpResult, l1)
				for i, c1 := range n1.Children {
					matrix[i] = make([]*CmpResult, l2)
					for j, c2 := range n2.Children {
						matrix[i][j] = MakePairsMock(a1, a2, c1, c2, pairs)
					}
				}

				first_used := make([]bool, l1)
				second_used := make([]bool, l2)

				max_i, max_j, max_res := bestOfSquareMatrix(matrix, first_used, second_used)
				for max_res != nil {
					//					max_rate := max_res.Result()
					first_used[max_i] = true
					second_used[max_j] = true
					pairs[id1][id2] = append(pairs[id1][id2], Coord{max_i, max_j})
					r.Append(*max_res)

					max_i, max_j, max_res = bestOfSquareMatrix(matrix, first_used, second_used)
				}
			}
			return r
		}
	}
	return nil
}
*/

func MergeShallow(n1 Node, n2 Node) *Node {
	if n1.Type == n2.Type {
		r := Node{
			Type: n1.Type,
			Attr: make([]html.Attribute, 0),
		}

		if n1.Data == n2.Data {
			r.Data = n1.Data
			for _, a1 := range n1.Attr {
				for _, a2 := range n2.Attr {
					if a1.Key == a2.Key && a1.Key != "class" {
						attr := html.Attribute{Key: a1.Key}
						if a1.Val == a2.Val {
							attr.Val = a1.Val
						}
						r.Attr = append(r.Attr, attr)
					}
				}
			}

			classes1 := n1.Classes()
			classes2 := n2.Classes()
			for _, c1 := range classes1 {
				for _, c2 := range classes2 {
					if c1 == c2 {
						r.AddClass(c1)
					}
				}
			}

		}

		if n1.DataArray != nil && n2.DataArray != nil {
			r.DataArray = append(n1.DataArray, n2.DataArray...)
		}

		return &r
	}
	return nil
}

/*
func DiffShallow(n1 Node, n2 Node) *Node {

}
*/

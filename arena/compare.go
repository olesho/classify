// classify project classify.go
package arena

import (
	"strings"

	"golang.org/x/net/html"
)

func bestOfRoundMatrix(matrix [][]int, checked []bool) (max_i int, max_j int, max_val int) {
	max_i = -1
	max_j = -1
	for i, _ := range matrix {
		if !checked[i] {
			for j := i + 1; j < len(matrix); j++ {
				if !checked[j] {
					if matrix[i][j] > max_val {
						max_val = matrix[i][j]
						max_i = i
						max_j = j
					}
				}
			}
		}
	}
	return
}

func bestOfSquareMatrix(matrix [][]*CmpResult, checked_i []bool, checked_j []bool) (max_i int, max_j int, max_val *CmpResult) {
	max_i = -1
	max_j = -1
	var max_sum float32
	for i, row := range matrix {
		if !checked_i[i] {
			for j, _ := range row {
				if !checked_j[j] && matrix[i][j] != nil {
					sum := matrix[i][j].Rate()
					if sum > max_sum {
						max_val = matrix[i][j]
						max_sum = sum
						max_i = i
						max_j = j
					}
				}
			}
		}
	}
	return
}

func CmpStrings(s1 string, s2 string) CmpResult {
	var r int
	l := len(s2)
	if len(s1) < len(s2) {
		l = len(s1)
	}

	for i := 0; i < l; i++ {
		if s1[i] == s2[i] {
			r++
		} else {
			break
		}
	}
	return CmpResult{r * 2, len(s1) + len(s2)}
}

func CmpDeepRate(a1 *Arena, a2 *Arena, id1, id2 int) *CmpResult {
	n1 := a1.Get(id1)
	n2 := a2.Get(id2)
	r := CmpShallow(n1, n2)

	if r != nil {
		child_r := CmpChildren(a1, a2, n1, n2)
		if child_r != nil {
			r.Append(*child_r)
		}
	}
	return r
}

type CmpResult struct {
	Sum   int
	Count int
}

func (r *CmpResult) Append(addition CmpResult) {
	r.Count += addition.Count
	r.Sum += addition.Sum
}

func (r *CmpResult) Rate() float32 {
	return float32(r.Sum) / float32(r.Count)
}

func CmpChildren(a1 *Arena, a2 *Arena, n1 *Node, n2 *Node) (r *CmpResult) {
	li := len(n1.Children)
	lj := len(n2.Children)

	if li == 0 && lj == 0 {
		return nil
	}

	r = &CmpResult{0, 0}

	if li > 0 && lj > 0 {
		matrix := make([][]*CmpResult, li)
		checked_i := make([]bool, li)
		checked_j := make([]bool, lj)
		for i := 0; i < li; i++ {
			matrix[i] = make([]*CmpResult, lj)
			for j := 0; j < lj; j++ {
				matrix[i][j] = CmpDeepRate(a1, a2, n1.Children[i], n2.Children[j])
			}
		}

		max_i, max_j, val := bestOfSquareMatrix(matrix, checked_i, checked_j)
		for val != nil {
			checked_i[max_i] = true
			checked_j[max_j] = true
			r.Count += val.Count
			r.Sum += val.Sum
			max_i, max_j, val = bestOfSquareMatrix(matrix, checked_i, checked_j)
		}

		for i := 0; i < li; i++ {
			if !checked_i[i] {
				r.Count++
			}
		}

		for j := 0; j < lj; j++ {
			if !checked_j[j] {
				r.Count++
			}
		}

		// double?
		//r.Rate *= 2

		return r
	}

	if li > 0 {
		r.Count = li
	}

	if lj > 0 {
		r.Count = lj
	}

	return r
}

func CmpShallow(n1, n2 *Node) *CmpResult {
	if n1.Type == n2.Type {
		if (n1.Type == html.TextNode) || (n1.Type == html.CommentNode) {
			r := CmpText(n1, n2)
			return &r
		}

		if n1.Data == n2.Data {
			r := &CmpResult{nodePoints, nodePoints}
			//return tagEqualityRate + (CmpAttr(n1.Attr, n2.Attr) * (1 - tagEqualityRate))
			attrResult := CmpAttr(n1.Attr, n2.Attr)
			if attrResult != nil {
				r.Append(*attrResult)
			}
			return r
		}
	}
	return nil
}

func CmpText(t1, t2 *Node) CmpResult {
	if t1.Data == t2.Data {
		return CmpResult{nodePoints + textPoints, nodePoints + textPoints}
	}
	return CmpResult{nodePoints, nodePoints + textPoints}
}

func CmpAttr(attr1, attr2 []html.Attribute) *CmpResult {
	if len(attr1) == 0 && len(attr2) == 0 {
		return nil
	}

	r := &CmpResult{0, 0}

	if len(attr1) > 0 && len(attr2) > 0 {
		//var r, r1, r2 int
		for _, a1 := range attr1 {
			for _, a2 := range attr2 {
				if a1.Key == a2.Key && a1.Key != "class" && a2.Key != "class" {
					r.Sum += attrKeyPoints
					valueRate := CmpStrings(a1.Val, a2.Val)
					valueRate.Count *= attrValPoints
					valueRate.Sum *= attrValPoints
					r.Append(valueRate)
				}
			}
		}

		classes1 := []string{}
		classes2 := []string{}
		for _, a1 := range attr1 {
			if a1.Key == "class" {
				classes1 = strings.Fields(a1.Val)
				r.Count += len(classes1) * classPoints
			} else {
				r.Count += attrKeyPoints + attrValPoints*len(a1.Val)
			}
		}
		for _, a2 := range attr2 {
			if a2.Key == "class" {
				classes2 = strings.Fields(a2.Val)
				r.Count += len(classes2) * classPoints
			} else {
				r.Count += attrKeyPoints + attrValPoints*len(a2.Val)
			}
		}

		for _, c1 := range classes1 {
			for _, c2 := range classes2 {
				if c1 == c2 {
					r.Sum += classPoints
				}
			}
		}

		//double?
		r.Sum *= 2

		return r
	}

	if len(attr1) > 0 {
		r.Count = len(attr1)
	}

	if len(attr2) > 0 {
		r.Count = len(attr2)
	}

	return r
}

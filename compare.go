// classify project classify.go
package classify

//	"fmt"

//	"golang.org/x/net/html"

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
	var max_sum float64
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

func (a *Arena) CmpDeepRate(id1, id2 int) *CmpResult {
	n1 := a.Get(id1)
	n2 := a.Get(id2)
	r := CmpShallow(n1, n2)

	if r != nil {
		child_r := a.CmpChildren(n1, n2)
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

func (r *CmpResult) Rate() float64 {
	return float64(r.Sum) / float64(r.Count)
}

func (a *Arena) CmpChildren(n1 Node, n2 Node) (r *CmpResult) {
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
				matrix[i][j] = a.CmpDeepRate(n1.Children[i], n2.Children[j])
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

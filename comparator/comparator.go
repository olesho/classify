package comparator

import (
	"sort"

	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
)

type Comparator interface {
	Cmp(idx1, idx2 int) float32
}

// DefaultComparator currently not used
type DefaultComparator struct {
	arena *arena.Arena
}

func NewDefaultComparator(a *arena.Arena) *DefaultComparator {
	return &DefaultComparator{a}
}

func (c *DefaultComparator) Cmp(idx1, idx2 int) float32 {
	n1, n2 := c.arena.Get(idx1), c.arena.Get(idx2)
	if c.cmpColumns(n1, n2) == 0 {
		return 0 // strict rule
	}
	val := c.cmpFull(n1, n2)
	return val * 2 / (GetVolume(n1) + GetVolume(n2))
}

func hasStr(s string, ss []string) bool {
	for _, n := range ss {
		if n == s {
			return true
		}
	}
	return false
}

func (s *DefaultComparator) cmpElements(n1, n2 *arena.Node) float32 {
	if n1.Type == n2.Type && n1.Type == html.TextNode {
		//return 0.1
		//return 1
		//return cmpStrings(n1.Data, n2.Data) * 2.0/float32(len(n1.Data) + len(n2.Data))
		//return cmpStrings(n1.Data, n2.Data)*10
	}
	if n1.Type == n2.Type {
		if n1.Data == n2.Data {
			var coincided float32 = 2.
			for _, attr1 := range n1.Attr {
				for _, attr2 := range n2.Attr {
					if attr1.Key == attr2.Key {
						coincided += 1
						if attr1.Key == "class" {
							classes1 := n1.Classes()
							classes2 := n2.Classes()
							for _, c1 := range classes1 {
								if hasStr(c1, classes2) {
									coincided += 1
								}
							}
						} else {
							coincided += cmpStrings(attr1.Val, attr2.Val)
						}
					}
				}
			}

			return coincided
		}
		return 0 // TODO to prevent any errors; but this should be changed for !COMMENT tag etc
	}
	return 0
}

func (s *DefaultComparator) cmpFull(n1, n2 *arena.Node) float32 {
	el := s.cmpElements(n1, n2)
	// strict tag names should coincide
	if el == 0 {
		return 0
	}
	ch := s.cmpChildren(n1, n2)
	return el + ch
}

func (s *DefaultComparator) cmpColumns(n1, n2 *arena.Node) float32 {
	chain1 := s.arena.Chain(n1.Id, 0)
	chain2 := s.arena.Chain(n2.Id, 0)
	size1 := len(chain1)
	size2 := len(chain2)
	if size1 != size2 {
		return 0
	}
	var r float32 = 0.
	for index := 1; (index < size1) && (index < size2); index++ {
		re := s.cmpElements(chain1[index], chain2[index])
		if re == 0 {
			return 0 // strict rule
		}
		r += re
	}
	return r / float32(size1)
}

func (s *DefaultComparator) cmpChildren(n1, n2 *arena.Node) float32 {
	size1, size2 := len(n1.Children), len(n2.Children)
	rating := make([]RateItem, size1*size2)
	for i1, idx1 := range n1.Children {
		for i2, idx2 := range n2.Children {
			idx := (i1+1)*(i2+1) - 1
			rc := s.cmpFull(s.arena.Get(idx1), s.arena.Get(idx2))
			rating[idx].Coincided = rc
			rating[idx].Index1 = i1
			rating[idx].Index2 = i2
		}
	}

	sort.Slice(rating, func(i, j int) bool {
		return rating[i].Coincided > rating[j].Coincided
	})

	flags1 := make([]bool, size1)
	flags2 := make([]bool, size2)
	count := 0
	smallerSize := 0
	if size1 < size2 {
		smallerSize = size1
	} else {
		smallerSize = size2
	}

	var coincided float32 = 0.
	for _, rate := range rating {
		if !flags1[rate.Index1] && !flags2[rate.Index2] {
			if rate.Coincided == 0 {
				break
			}
			coincided += rate.Coincided
			flags1[rate.Index1] = true
			flags2[rate.Index2] = true
			count++
			if count == smallerSize {
				break
			}
		}
	}

	return coincided
}

type RateItem struct {
	Coincided float32
	Index1    int
	Index2    int
}

type Ratio struct {
	Num float32
	Den float32
}

type Result struct {
	Coincided float32
	Total     float32
}

func cmpStrings(s1 string, s2 string) float32 {
	if len(s1) == 0 && len(s2) == 0 {
		return 0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}

	var coincided float32
	l := len(s2)
	if len(s1) < len(s2) {
		l = len(s1)
	}

	for i := 0; i < l; i++ {
		if s1[i] == s2[i] {
			coincided += 1.
		} else if isDigitChar(s1[i]) && isDigitChar(s2[i]) { // both are digits
			coincided += 0.8
		} else if isUpperChar(s1[i]) && isUpperChar(s2[i]) { // both are upper characters
			coincided += 0.5
		} else if isLowerChar(s1[i]) && isLowerChar(s2[i]) { // both are lower characters
			coincided += 0.5
		} else if (isUpperChar(s1[i]) || isLowerChar(s1[i])) && (isUpperChar(s2[i]) || isLowerChar(s2[i])) { // both are characters
			coincided += 0.2
		} else {
			break
		}
	}
	return coincided * 2 / float32(len(s1)+len(s2)) + 1 // +1 if both strings are nonempty
}

func isLowerChar(c byte) bool {
	return c > 96 && c < 123
}

func isUpperChar(c byte) bool {
	return c > 32 && c < 91
}

func isDigitChar(c byte) bool {
	return c > 47 && c < 58
}

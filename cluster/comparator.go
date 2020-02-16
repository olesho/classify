package cluster

import (
	classify "github.com/olesho/class"
	"sort"
)

type Comparator interface {
	Cmp(n1, n2 *classify.Node) float64
}

type DefaultComparator struct {
	arena *classify.Arena
}

func NewDefaultComparator(a *classify.Arena) *DefaultComparator {
	return &DefaultComparator{ a }
}

func (c *DefaultComparator) Cmp(n1, n2 *classify.Node) float64 {
	ce := c.cmpElements(n1, n2)
	if ce.Similarity == 0 {
		return 0
	}

	rc := c.cmpColumns(n1, n2)
	if rc == 0 {
		return 0 // strict rule
	}

	cr := c.cmpChildren(n1, n2)
	//return rc * 0.1 + ce.Similarity + cr.Similarity
	return cr.Coincided
}

func (s *DefaultComparator) cmpElements(n1, n2 *classify.Node) Result {
	denominator := tokenVolume(n1) + tokenVolume(n2)
	if n1.Type == n2.Type {
		if n1.Data == n2.Data {
			total := 2.
			for _, attr1 := range n1.Attr {
				for _, attr2 := range n2.Attr {
					if attr1.Key == attr2.Key {
						total += 1
						total += cmpStrings(attr1.Val, attr2.Val).Similarity
					}
				}
			}
			return Result{total, total*2 / denominator}
		}
		return Result{0,0} // TODO to prevent any errors; but this should be changed for !COMMENT tag etc
	}
	return Result{0,0}
}

func (s *DefaultComparator) cmpColumns(n1, n2 *classify.Node)  float64 {
	chain1 := s.arena.Chain(n1.Id, 0)
	chain2 := s.arena.Chain(n2.Id, 0)
	size1 := len(chain1)
	size2 := len(chain2)
	if size1 != size2 {
		return 0
	}
	r := 0.
	for index := 1; (index < size1) && (index < size2); index++ {
		re := s.cmpElements(chain1[index], chain2[index])
		if re.Similarity == 0 {
			return 0 // strict rule
		}
		r += re.Similarity
	}
	return r/float64(size1)
}


func (s *DefaultComparator) cmpChildren(n1, n2 *classify.Node) Result {
	size1, size2 := len(n1.Children), len(n2.Children)
	rating := make([]RateItem, size1*size2)
	for i1, idx1 := range n1.Children {
		for i2, idx2 := range n2.Children {
			idx := (i1+1)*(i2+1) - 1
			rc := s.cmpElements(s.arena.Get(idx1), s.arena.Get(idx2))
			rating[idx].Similarity, rating[idx].Coincided = rc.Similarity, rc.Coincided
			rating[idx].Index1 = i1
			rating[idx].Index2 = i2
		}
	}

	sort.Slice(rating, func(i, j int) bool {
		return rating[i].Similarity > rating[j].Similarity
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

	result := Ratio{}
	for _, rate := range rating {
		if !flags1[rate.Index1] && !flags2[rate.Index2] {
			result.Num += rate.Similarity
			result.Den += 2
			flags1[rate.Index1] = true
			flags2[rate.Index2] = true
			count++
			if count == smallerSize {
				break
			}
		}
	}

	return Result{
		Coincided: result.Num,
		Similarity: result.Num/result.Den,
	}
}

type RateItem struct {
	Result
	Index1 int
	Index2 int
}

type Ratio struct {
	Num float64
	Den float64
}

type Result struct {
	Coincided float64
	Similarity float64
}


// this counts inform
func tokenVolume(n *classify.Node) float64 {
	volume := 1. // has Type
	if len(n.Data) > 1 { // has Data
		volume += 1
	}
	for _, attr := range n.Attr { // has Attributes
		if len(attr.Key) > 0 {
			volume += 1
		}
		volume += float64(len(attr.Val))
		// TODO class and id might be treated differently
	}
	return volume
}

func cmpStrings(s1 string, s2 string) Result {
	if len(s1) == 0 && len(s2) == 0 {
		return Result{0,0}
	}
	if len(s1) == 0 || len(s2) == 0 {
		return Result{0,0}
	}

	var coincided float64
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
	return Result{
		Coincided: coincided,
		Similarity: coincided * 2 / float64(len(s1) + len(s2)),
	}
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


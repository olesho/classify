package stream

import (
	"fmt"
	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
	"sort"
	"strings"
)

type Cluster struct {
	Indexes []int
	Rate    float32
}

func (mtx *Mtx) FindIndexes(indexes []int) []int {
	r := make([]int, len(indexes))
	for i, idx := range indexes {
		r[i] = mtx.Indexes[idx]
	}
	return r
}

func (s *Storage) toTable(c *Cluster) Table {
	if c.hasIndex(1217) {
		fmt.Println()
	}

	templateArena := s.MergeAll(c, c.Indexes)
	result := Table{
		Arena:         s.Arena,
		TemplateArena: templateArena,
		Members:       make([]*arena.Node, len(c.Indexes)),
		Rate:          c.Rate,
	}

	result.Fields = result.WholesomeGroupFields()
	for i, memberIdx := range c.Indexes {
		result.Members[i] = s.Arena.Get(memberIdx)
	}

	return result
}

// MergeAll merges all nodes with indexes into single template producing new arena
func (s *Storage) MergeAll(c *Cluster, indexes []int) *arena.Arena {
	rootID := indexes[0]
	templateArena := initClone(s.Arena, rootID)
	for _, nextID := range indexes[1:] {
		templateArena.List[1].Ext.(*Additional).AppendGroupId(nextID)
		s.MergeIntoTemplate(c, templateArena, nextID, 1)
	}
	return templateArena
}

func initClone(arena *arena.Arena, root int) *arena.Arena {
	res := arena.CloneBranch(root)
	Init(res)
	for i := range res.List[1:] {
		res.List[1:][i].Ext.(*Additional).AppendGroupId(i + root)
	}
	return res
}


func mergeIntoTemplateAttrs(node1, node2 *arena.Node) []html.Attribute {
	mergedAttrs := []html.Attribute{}
	for _, attr1 := range node1.Attr {
		for _, attr2 := range node2.Attr {
			if attr1.Key == attr2.Key {
				mergedAttr := html.Attribute{
					Key: attr1.Key,
				}
				if mergedAttr.Key == "class" {
					mergedClasses := []string{}
					for _, class1 := range node1.Classes() {
						for _, class2 := range node2.Classes() {
							if class1 == class2 {
								mergedClasses = append(mergedClasses, class1)
							}
						}
					}
					mergedAttr.Val = strings.Join(mergedClasses, " ")
				} else {
					if attr1.Val == attr2.Val {
						mergedAttr.Val = attr1.Val
					}
				}
				mergedAttrs = append(mergedAttrs, mergedAttr)
			}
		}
	}
	return mergedAttrs
}

func (s *Storage) MergeIntoTemplate(c *Cluster, templateArena *arena.Arena, mainIdx, templateIdx int) {
	n1 := s.Arena.Get(mainIdx)
	templateNode := templateArena.Get(templateIdx)
	n2 := s.Arena.Get(templateNode.Ext.(*Additional).GroupIds[0])

	// merge attributes
	templateArena.List[templateIdx].Attr = mergeIntoTemplateAttrs(templateNode, n1)

	// merge text node
	if templateNode.Type == html.TextNode && n1.Type == html.TextNode {
		if templateNode.Data != n1.Data {
			templateNode.Data = "#text"
		}
	}

	// merge children
	size1, size2 := len(n1.Children), len(n2.Children)
	ratingMatrixSize := size1 * size2
	if ratingMatrixSize > 0 {
		rating := make([]mergeItem, ratingMatrixSize)
		for i1, idx1 := range n1.Children {
			for i2, idx2 := range n2.Children {
				idx := (i1+1)*(i2+1) - 1
				rating[idx].Similarity = s.Find(idx1, idx2)
				rating[idx].Index1 = i1
				rating[idx].Index2 = i2
				rating[idx].TemplateIndex = templateNode.Children[i2]
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

		for _, rate := range rating { // TODO add exit if rate.Similarity == 0 !!!
			if !flags1[rate.Index1] && !flags2[rate.Index2] {
				if rate.Similarity == 0 {
					break
				}
				templateChildNode := templateArena.Get(templateNode.Children[rate.Index2])
				idx := n1.Children[rate.Index1]
				templateChildNode.Ext.(*Additional).AppendGroupId(idx)
				s.MergeIntoTemplate(c, templateArena, idx, rate.TemplateIndex)
				flags1[rate.Index1] = true
				flags2[rate.Index2] = true
				count++
				if count == smallerSize {
					break
				}
			}
		}
	}
}

type mergeItem struct {
	Similarity    float32
	Index1        int
	Index2        int
	TemplateIndex int
}

func (c *Cluster) tryAdd(candidateRate float32, candidateIndex int) bool {
	if c.Volume() < candidateRate*float32(len(c.Indexes)+1) {
		c.Rate = candidateRate
		c.Indexes = append(c.Indexes, candidateIndex)
		return true
	}
	return false
}

func (c *Cluster) add(candidateRate float32, candidateIndex int) *Cluster {
	return &Cluster{
		Indexes: append(c.Indexes, candidateIndex),
		Rate: candidateRate,
	}
}

func (c *Cluster) clone() *Cluster {
	cc := &Cluster{
		Indexes: make([]int, len(c.Indexes)),
		Rate: c.Rate,
	}
	copy(cc.Indexes, c.Indexes)
	return cc
}

func (c *Cluster) Volume() float32 {
	return float32(len(c.Indexes)) * c.Rate
}
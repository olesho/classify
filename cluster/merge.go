package cluster

import (
	"github.com/olesho/classify"
	"golang.org/x/net/html"
	"sort"
	"strings"
)

func MergeAll(arena *classify.Arena, matrix *RateMatrix, indexes []int) *classify.Arena {
	rootId := indexes[0]
	templateArena := initClone(arena, matrix, rootId)
	for _, nextId := range indexes[1:] {
		templateArena.List[0].Ext.(*Additional).AppendGroupId(nextId)
		mergeIntoTemplate(arena, templateArena, nextId, 0, matrix)
	}
	return templateArena
}

func initClone(arena *classify.Arena, matrix *RateMatrix, root int) *classify.Arena {
	res := arena.CloneBranch(root)
	Init(res)
	for i := range res.List {
		res.List[i].Ext.(*Additional).AppendGroupId(i + root)
	}
	return res
}

type mergeItem struct {
	Similarity    float64
	Index1        int
	Index2        int
	TemplateIndex int
}

func mergeIntoTemplateAttrs(node1, node2 *classify.Node) []html.Attribute {
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

func mergeIntoTemplate(mainArena, templateArena *classify.Arena, mainIdx, templateIdx int, matrix *RateMatrix) {
	n1 := mainArena.Get(mainIdx)
	templateNode := templateArena.Get(templateIdx)
	n2 := mainArena.Get(templateNode.Ext.(*Additional).GroupIds[0])

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
	rating := make([]mergeItem, size1*size2)
	for i1, idx1 := range n1.Children {
		for i2, idx2 := range n2.Children {
			idx := (i1+1)*(i2+1) - 1
			rating[idx].Similarity = matrix.Cmp(idx1, idx2)
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

	for _, rate := range rating { // TODO add exit if rate.Similariry == 0 !!!
		if !flags1[rate.Index1] && !flags2[rate.Index2] {
			if rate.Similarity == 0 {
				break
			}
			templateChildNode := templateArena.Get(templateNode.Children[rate.Index2])
			idx := n1.Children[rate.Index1]
			templateChildNode.Ext.(*Additional).AppendGroupId(idx)
			mergeIntoTemplate(mainArena, templateArena, idx, rate.TemplateIndex, matrix)
			flags1[rate.Index1] = true
			flags2[rate.Index2] = true
			count++
			if count == smallerSize {
				break
			}
		}
	}
}

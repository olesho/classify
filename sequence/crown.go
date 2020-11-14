package sequence

import (
	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
	"sort"
	"strings"
)

type CrownCluster struct {
	indexes []int
	stem *StemCluster
	rate float32
}

func (c *CrownCluster) Has(index int) bool {
	for _, nextIndex := range c.indexes {
		if nextIndex == index {
			return true
		}
	}
	return false
}

func (c *CrownCluster) HasUnresolved(index int) bool {
	for _, nextIndex := range c.indexes {
		if c.stem.indexes[nextIndex] == index {
			return true
		}
	}
	return false
}

func (c *CrownCluster) resolveIndexes() {
	for i, index := range c.indexes {
		c.indexes[i] = c.stem.indexes[index]
	}
}


func (c *CrownCluster) String() string {
	r := make([]string, len(c.indexes))
	for i, index := range c.indexes {
		r[i] = c.stem.root.Arena.StringifyNode(index)
	}
	return strings.Join(r, " ")
}

func (c *CrownCluster) Volume() float32 {
	return float32(len(c.indexes)) * c.rate
}

func (c *CrownCluster) attachDelta(newRate float32) float32 {
	return newRate * float32(len(c.indexes)+1) - float32(len(c.indexes)) * c.rate
}

func (c *CrownCluster) attach(rate float32, localIndex int)  {
	if len(c.indexes) == 2 || rate < c.rate {
		c.rate = rate
	}
	c.indexes = append(c.indexes, localIndex)
}

func (c *CrownCluster) detach(rate float32, localIndex int)  {
	c.rate = rate
	c.indexes = append(c.indexes[:localIndex], c.indexes[localIndex+1:]...)
}

func (c *CrownCluster) detachRate() (crownIndex int, delta, nextRate float32) {
	if len(c.indexes) > 2 {
		var lowestIndex = -1
		var lowestSum float32

		for n, i := range c.indexes {
			var sum float32
			for _, j := range c.indexes {
				if i != j {
					sum += c.stem.Get(i, j)
				}
			}
			if lowestSum == 0 {
				lowestSum = sum
				lowestIndex = n
			} else {
				if sum < lowestSum {
					lowestSum = sum
					lowestIndex = n
				}
			}
		}

		var secondLowestRate float32
		for n, i := range c.indexes {
			if lowestIndex != n {
				for m, j := range c.indexes {
					if lowestIndex != m && i != j {
						val := c.stem.Get(i, j)
						if secondLowestRate == 0 {
							secondLowestRate = val
						} else {
							if val < secondLowestRate{
								secondLowestRate = val
							}
						}
					}
				}
			}
		}

		return lowestIndex, secondLowestRate*float32(len(c.indexes)-1) - c.rate * float32(len(c.indexes)), secondLowestRate
	}
	return -1, 0, 0
}

type pair struct {
	index int
	val float32
}

func (c *CrownCluster) extend() {
	var pairs = make([]pair, 0, len(c.stem.indexes))
	for i := range c.stem.indexes {
		if !c.Has(i) {
			pairs = append(pairs, pair{index: i, val: c.attachRate(i)})
		}
	}

	sort.Slice(pairs, func (i, j int) bool {
		return pairs[i].val > pairs[j].val
	})

	for _, pair := range pairs {
		delta := c.attachDelta(pair.val)
		if delta > 0 {
			c.attach(pair.val, pair.index)
		} else {
			break
		}
	}
	sort.Ints(c.indexes)
}

func (c *CrownCluster) attachRate(stemIndex int) float32 {
	var lowestVal float32
	var i int
	if len(c.indexes) > c.stem.root.limit {
		i = len(c.indexes) - c.stem.root.limit
	}

	for ; i < len(c.indexes); i++ {
		v := c.stem.Get(stemIndex, c.indexes[i])
		if v > 0 {
			if lowestVal == 0 {
				lowestVal = v
				continue
			}
			if v < lowestVal {
				lowestVal = v
			}
		}
		//} else {
		// removed because of limit; only last [limit] values will be > 0
		//return 0
		//}
	}
	return lowestVal
}


func (c *CrownCluster) toTable() Table {
	templateArena := c.MergeAll(c.indexes)
	result := Table{
		Arena:         c.stem.root.Arena,
		TemplateArena: templateArena,
		Members:       make([]*arena.Node, len(c.indexes)),
		Rate:          c.rate,
	}

	result.Fields = result.WholesomeGroupFields()
	for i, memberIdx := range c.indexes {
		result.Members[i] = c.stem.root.Arena.Get(memberIdx)
	}

	return result
}

// MergeAll merges all nodes with indexes into single template producing new arena
func (c *CrownCluster) MergeAll(indexes []int) *arena.Arena {
	templateArena := initClone(c.stem.root.Arena, indexes[0])
	for _, nextID := range indexes[1:] {
		templateArena.List[1].Ext.(*Additional).AppendGroupId(nextID)
		c.MergeIntoTemplate(templateArena, nextID, 1)
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

func mergeIntoTemplateAttrs(n1, n2 *arena.Node) []html.Attribute {
	var mergedAttrs []html.Attribute
	for _, attr1 := range n1.Attr {
		for _, attr2 := range n2.Attr {
			if attr1.Key == attr2.Key {
				mergedAttr := html.Attribute{
					Key: attr1.Key,
				}
				if mergedAttr.Key == "class" {
					mergedClasses := []string{}
					for _, class1 := range n1.Classes() {
						for _, class2 := range n2.Classes() {
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

func (c *CrownCluster) MergeIntoTemplate(templateArena *arena.Arena, mainIdx, templateIdx int) {
	n1 := c.stem.root.Arena.Get(mainIdx)
	templateNode := templateArena.Get(templateIdx)

	// merge attributes
	templateArena.List[templateIdx].Attr = mergeIntoTemplateAttrs(templateNode, n1)

	// merge text node
	if templateNode.Type == html.TextNode && n1.Type == html.TextNode {
		if templateNode.Data != n1.Data {
			templateNode.Data = "#text"
		}
	}

	// merge children
	size1, size2 := len(n1.Children), len(templateNode.Children)
	ratingMatrixSize := size1 * size2
	if ratingMatrixSize > 0 {
		rating := make([]mergeItem, 0, ratingMatrixSize)
		for i2, idx2 := range templateNode.Children {
			templateChildNodeGroupIDs := templateArena.Get(idx2).Ext.(*Additional).GroupIds
			for i1, idx1 := range n1.Children {
				//idx := (i1+1)*(i2+1) - 1

				var sum float32
				cnt := 0
				for _, id := range templateChildNodeGroupIDs {
					if val := c.stem.root.FindCrown(idx1, id); val > 0  {
						sum += val
						cnt++
					}
				}

				item := mergeItem{
					Index1: i1,
					Index2: i2,
				}

				if cnt > 0 {
					item.Similarity = sum / float32(cnt)
				}
				rating = append(rating, item)
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
				c.MergeIntoTemplate(templateArena, idx, templateNode.Children[rate.Index2])
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
}

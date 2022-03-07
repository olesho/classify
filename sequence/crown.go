package sequence

import (
	"github.com/olesho/classify/arena"
	"golang.org/x/net/html"
	"sort"
	"strings"
)

type CrownItem struct {
	Index int
	ResolvedIndex int
	ValueSum float32
}

type CrownCluster struct {
	items []CrownItem
	stem    *StemCluster
	rate    float32
}

func (c *CrownCluster) Has(index int) bool {
	for _, item := range c.items {
		if item.Index == index {
			return true
		}
	}
	return false
}

func (c *CrownCluster) HasResolved(index int) bool {
	for _, item := range c.items {
		if item.ResolvedIndex == index {
			return true
		}
	}
	return false
}


func (c *CrownCluster) resolveIndexes() {
	for i := range c.items {
		//c.indexes[i] = c.stem.indexes[index]
		c.items[i].ResolvedIndex = c.stem.indexes[c.items[i].Index]
	}
}

func (c *CrownCluster) String() string {
	r := make([]string, len(c.items))
	for i := range c.items {
		r[i] = c.stem.root.Arena.StringifyNode(c.items[i].Index)
	}
	return strings.Join(r, " ")
}

func (c *CrownCluster) Volume() float32 {
	return float32(len(c.items)) * c.rate
}

func (c *CrownCluster) ExpandBest(low float32, localIndex int) bool {
	volume, nextVolume := float32(len(c.items))*c.rate, low*float32(len(c.items)+1)
	if volume < nextVolume {
		if c.rate == 0 || low < c.rate {
			c.rate = low
		}

		var start int
		if len(c.items) > c.stem.root.limit {
			start = len(c.items) - c.stem.root.limit
		}

		var newItemValuesSum float32
		for i := start; i < len(c.items); i++ {
			v := c.stem.Get(localIndex, c.items[i].Index)
			c.items[i].ValueSum += v
			newItemValuesSum += v
		}

		c.items = append(c.items, CrownItem{
			Index: localIndex,
			ValueSum: newItemValuesSum,
		})

		return true
	}
	return false
}

func (c *CrownCluster) SqueezeWorst() (squeezedIndex int) {
	squeezedIndex = -1
	if len(c.items) > 2 {
		minIdx := 0
		minAvg := c.items[0].ValueSum / float32(len(c.items))
		for i := range c.items[1:] {
			v := c.items[i+1].ValueSum  / float32(len(c.items))
			if v < minAvg {
				minIdx = i+1
				minAvg = v
			}
		}

		var newLow float32
		for i := range c.items {
			if i != minIdx {
				for j := i+1; j < len(c.items); j++ {
					if j != minIdx && i != j {
						nextRate := c.stem.Get(c.items[i].Index, c.items[j].Index)
						if newLow == 0 || nextRate < newLow {
							newLow = nextRate
						}
					}
				}
			}
		}

		if newLow * float32(len(c.items)-1) > c.Volume() {
			oldIdx := c.items[minIdx].Index

			c.items = append(c.items[:minIdx], c.items[minIdx+1:]...)

			for i := range c.items {
				squeezedVal := c.stem.Get(oldIdx, c.items[i].Index)
				c.items[i].ValueSum -= squeezedVal
			}

			c.rate = newLow
			return minIdx
		}
	}
	return
}

func (c *CrownCluster) Rate(stemIndex int) (low float32, sum float32) {
	var start int
	if len(c.items) > c.stem.root.limit {
		start = len(c.items) - c.stem.root.limit
	}

	sum = 0
	for i := start; i < len(c.items); i++ {
		v := c.stem.Get(stemIndex, c.items[i].Index)
		sum += v
		if v > 0 {
			if low == 0 {
				low = v
				continue
			}
			if v < low {
				low = v
			}
		}
		//} else {
		// removed because of limit; only last [limit] values will be > 0
		//return 0
		//}
	}
	return low, sum
}

func (c *CrownCluster) toTable() Table {
	resolvedIndexes := make([]int, len(c.items))
	for i := range c.items { resolvedIndexes[i] = c.items[i].ResolvedIndex }
	templateArena := c.MergeAll(resolvedIndexes)
	result := Table{
		Arena:         c.stem.root.Arena,
		TemplateArena: templateArena,
		Members:       make([]*arena.Node, len(c.items)),
		Rate:          c.rate,
	}

	result.FieldSets = result.WholesomeGroupFields()
	for i, memberIdx := range resolvedIndexes {
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
					if val := c.stem.root.FindCrown(idx1, id); val > 0 {
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
	Similarity float32
	Index1     int
	Index2     int
}

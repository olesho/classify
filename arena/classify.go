// classify project classify.go
package arena

//type Classificator interface {
//	Run()
//	Bags() Bags
//}
//
///*
//func GenerateFields(a *Arena, bag Bag) [][]string {
//	DiffShallow()...
//}
//*/
//
//func GeneratePattern(a *Arena, bag Bag) *Arena {
//	bagSize := len(bag.Content)
//	if bagSize > 1 {
//		pattern := a.Clone(bag.Content[0])
//		for n := 1; n < bagSize; n++ {
//			next := a.Clone(bag.Content[n])
//			pattern = Merge(pattern, next, 0, 0)
//		}
//		return pattern
//	}
//	return nil
//}
//
//func GeneratePath(a *Arena, bag Bag) []string {
//	return generatePath(a, bag.Content, make([]string, 0))
//}
//
//func generatePath(a *Arena, nodes []int, currentPath []string) []string {
//	nodesSize := len(nodes)
//	if nodesSize > 0 {
//		parents := make([]int, nodesSize)
//		patternNode := a.Get(nodes[0])
//		pattern := &patternNode
//		parents[0] = a.Get(nodes[0]).Parent
//		for i := 1; i < nodesSize; i++ {
//			pattern = MergeShallow(*pattern, a.Get(nodes[i]))
//			if pattern.Data == "" {
//				return currentPath
//			}
//			parents[i] = a.Get(nodes[i]).Parent
//		}
//		path := append(currentPath, pattern.String())
//		for _, p := range parents {
//			if p == 0 {
//				return path
//			}
//		}
//		return generatePath(a, parents, path)
//	}
//	return nil
//}
//
//func IsParent(a *Arena, parent, child Bag) bool {
//	if len(parent.Content) == len(child.Content) {
//		has := 0
//		for _, c := range child.Content {
//			for _, p := range parent.Content {
//				if a.HasParent(c, p) {
//					has++
//				}
//			}
//		}
//		if has == len(parent.Content) {
//			return true
//		}
//	}
//	return false
//}

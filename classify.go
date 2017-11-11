// classify project classify.go
package classify

type Classificator interface {
	Run()
	Bags() Bags
}

func Pattern(a *Arena, bag Bag) *Arena {
	bagSize := len(bag.Content)
	if bagSize > 1 {
		pattern := a.Clone(bag.Content[0])
		for n := 1; n < bagSize; n++ {
			next := a.Clone(bag.Content[n])
			pattern = Merge(pattern, next, 0, 0)
		}
		return pattern
	}
	return nil
}

func IsParent(a *Arena, parent, child Bag) bool {
	if len(parent.Content) == len(child.Content) {
		has := 0
		for _, c := range child.Content {
			for _, p := range parent.Content {
				if a.HasParent(c, p) {
					has++
				}
			}
		}
		if has == len(parent.Content) {
			return true
		}
	}
	return false
}

package sequence

import "github.com/olesho/classify/arena"

type Additional struct {
//	Volume   float32
//	GroupIds []int
	LastDescendant int
}

func initExt(a *arena.Arena) {
	for i := range a.List {
		a.List[i].Ext = &Additional{}
	}
}

func initNode(a *arena.Arena, idx int) int {
	n := a.Get(idx)
	if len(n.Children) > 0 {
		for i := 0; i < len(n.Children)-1; i++ {
			initNode(a, n.Children[i])
		}
		n.Ext.(*Additional).LastDescendant = initNode(a, n.Children[len(n.Children)-1])
	} else {
		n.Ext.(*Additional).LastDescendant = idx
	}
	return n.Ext.(*Additional).LastDescendant
}

func Init(a *arena.Arena) {
	initExt(a)
	initNode(a, 0)
}

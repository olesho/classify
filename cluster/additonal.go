package cluster

import "github.com/olesho/classify/arena"

type Additional struct {
	Volume   float64
	GroupIds []int
}

func (a *Additional) AppendGroupId(id int) {
	a.GroupIds = append(a.GroupIds, id)
}

func (a *Additional) AppendVolume(v float64) {
	a.Volume += v
}

func GetVolume(n *arena.Node) float64 {
	return n.Ext.(*Additional).Volume
}

func Init(a *arena.Arena) {
	initExt(a)
	initVolume(a)
}

func initExt(a *arena.Arena) {
	for i := range a.List {
		a.List[i].Ext = &Additional{}
	}
}

func initVolume(a *arena.Arena) {
	for _, el := range a.List {
		el.Ext.(*Additional).AppendVolume(tokenVolume(el))
	}
	for id := len(a.List) - 1; id > -1; id-- {
		el := a.List[id]
		for _, childIdx := range el.Children {
			el.Ext.(*Additional).AppendVolume(GetVolume(a.List[childIdx]))
		}
	}
}

func tokenVolume(n *arena.Node) float64 {
	volume := .5         // has Type
	if len(n.Data) > 1 { // has Data
		volume += .5
	}
	for _, attr := range n.Attr { // has Attributes
		if len(attr.Key) > 0 {
			volume += 1
		}
		if attr.Key == "class" {
			volume += float64(len(n.Classes()))
		} else {
			volume += float64(len(attr.Val))
		}
	}
	return volume
}

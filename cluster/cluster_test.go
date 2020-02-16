package cluster

import (
	"bufio"
	"fmt"
	classify "github.com/olesho/class"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"os"
	"testing"
)

func TestYcomb(t *testing.T) {
	a := assert.New(t)

	f, _ := os.Open("ycomb.html")
	defer f.Close()
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	a.NoError(err)

	arena := classify.NewArena(*n)
	m := Clusterize(arena)
	rank := m.Rank(0)
	for _, row := range rank.Matrix {
		for _, n := range row {
			str, _ := arena.RenderString(n.Id)
			fmt.Println(str)
		}
		fmt.Println("__________________________________________________________________________________________________")
	}
}

//
//func isin(candidates []Cell, id int) bool {
//	for _, c := range candidates {
//		if c.Index == id {
//			return true
//		}
//	}
//	return false
//}
//
//func cnt(nodes []*classify.Node, id int) int {
//	c := 0
//	for _, n := range nodes {
//		if n.Id == id {
//			c++
//		}
//	}
//	return c
//}
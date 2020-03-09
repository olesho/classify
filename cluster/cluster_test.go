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

	//f, _ := os.Open("examples/ycomb.html")
	//f, _ := os.Open("examples/hackernoon.html")
	//f, _ := os.Open("examples/pravda.html")
	f, _ := os.Open("examples/bbc.html")
	defer f.Close()
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	a.NoError(err)

	arena := classify.NewArena(*n)
	m := Clusterize(arena)
	rank := m.Rank(0)
	for _, row := range rank.Nonuniform().Matrix {
		for _, n := range row {
			//str, _ := arena.RenderString(n.Id)
			str := arena.StringifyInformation(n.Id)
			fmt.Println(str)
		}
		fmt.Println("__________________________________________________________________________________________________")
	}
}

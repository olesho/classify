package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/olesho/classify"
	"github.com/olesho/classify/bags"
	"golang.org/x/net/html"
)

// func TestMain(t *testing.T) {
// 	f, _ := os.Open("fb.html")

// 	reader := bufio.NewReader(f)
// 	n, err := html.Parse(reader)
// 	if err != nil {
// 		panic(err)
// 	}

// 	arena := classify.NewArena(*n)
// 	fmt.Println(bags.ExtendedComparator(arena, arena.Get(6221), arena.Get(6561)))
// }

func TestMain(t *testing.T) {
	f, _ := os.Open("fb.html")

	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	if err != nil {
		panic(err)
	}

	arena := classify.NewArena(*n)
	fmt.Println(bags.ExtendedComparator(arena, arena.Get(315), arena.Get(308)))
	//fmt.Println(bags.ExtendedComparator(arena, arena.Get(308), arena.Get(308)))

	// matrix, err := bags.Parse(arena)
	// if err != nil {
	// 	panic(err)
	// }

	//idxs := []int{}
	// for _, row := range matrix[0] {
	// 	for _, n := range row {
	// 		//idxs = append(idxs, n.Id)
	// 		fmt.Println(arena.StringifyNode(n.Id))
	// 	}
	// }

	// for _, idx := range idxs {
	// 	fmt.Println(bags.ExtendedComparator(arena, arena.Get(idx), arena.Get(idxs[0])))
	// }
}

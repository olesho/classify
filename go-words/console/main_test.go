package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/olesho/classify/bags"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

func TestMain(t *testing.T) {
	f, _ := os.Open("fb.html")

	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	if err != nil {
		panic(err)
	}

	arena := classify.NewArena(*n)

	//fmt.Println(bags.ExtendedComparator(arena, arena.Get(97), arena.Get(16905)))

	fmt.Println(bags.ExtendedComparator(arena, arena.Get(6221), arena.Get(6561)))
}

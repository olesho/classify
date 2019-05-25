package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/olesho/classify"
	"github.com/olesho/classify/bags"
	"golang.org/x/net/html"
)

func main() {
	rank := flag.Int("r", 0, "Rank")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	n, err := html.Parse(reader)
	if err != nil {
		panic(err)
	}

	arena := classify.NewArena(*n)

	matrix, err := bags.Parse(arena)
	if err != nil {
		panic(err)
	}

	if len(matrix) >= *rank+1 {

		// transform
		fmt.Println("WIDTH:", len(matrix[*rank][0]))
		fmt.Println("LENGTH:", len(matrix[*rank]))
		fmt.Println("===============================================")
		for _, row := range matrix[*rank] {
			for _, n := range row {
				fmt.Println(arena.StringifyInformation(n.Id))
				fmt.Println()
			}
			fmt.Println("===============================================")
		}
	}

	// YCOMBINATOR PROBLEM !!!

	//for _, n := range nodes {
	//	//fmt.Println(arena.RenderString(n.Id))
	//
	//	fmt.Println(arena.StringifyInformation(n.Id))
	//	fmt.Println()
	//	fmt.Println("===============================================")
	//}
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	classify "github.com/olesho/class"
	"github.com/olesho/class/bags"
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

	m, err := bags.Parse(arena)
	if err != nil {
		panic(err)
	}

	if rank := m.Rank(*rank); rank != nil {

		// transform
		//fmt.Println("WIDTH:", len(rank.Matrix[0]))
		//fmt.Println("LENGTH:", len(rank.Matrix))
		//fmt.Println("===============================================")
		for _, row := range rank.Nonuniform().Matrix {
			//fmt.Println(len(row))
			for _, n := range row {
				//fmt.Println(arena.StringifyNode(n.Id))
				//fmt.Println(arena.StringifyInformation(n.Id))
				fmt.Print(`"`+strings.Replace(strings.TrimSpace(arena.StringifyInformation(n.Id)), "\n", "", -1)+`"`, ",")
				//fmt.Println(arena.StringifyWithChildren(n.Id))
				//fmt.Println("-----------------------------------------------")
			}
			//fmt.Println("===============================================")
			fmt.Println()
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

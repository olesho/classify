package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/olesho/classify/arena"
	"github.com/olesho/classify/cluster"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var defaultArena *arena.Arena
var matrix *cluster.Matrix

func FuncClear(command string) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
		return
	}
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func renderOutput(groupIdx int, render func(itemID int) []string) {
	cols, _ := consolesize.GetConsoleSize()
	bw := NewBlockWriter(cols+1, 1, 1)
	bw.Open(prompt.Red, prompt.White, fmt.Sprintf("Total groups:%v", len(matrix.Matrix[groupIdx].Matrix)))
	for i, series := range matrix.Matrix[groupIdx].Matrix {
		bw.Open(prompt.Brown, prompt.White, fmt.Sprintf("Group %v", i+1))
		for itemIndex, item := range series {
			subItems := render(item.Id)
			bw.Open(prompt.Cyan, prompt.White, fmt.Sprintf("Item %v", itemIndex+1))
			for _, subItem := range subItems {
				subItem = strings.ReplaceAll(subItem, "\n", " ")
				bw.WriteText(prompt.Black, prompt.White, false, subItem)
			}
			bw.Close()
		}
		bw.Close()
	}
	bw.Close()
	return
}

func FuncShow(command string) {
	if matrix != nil {
		parts := showRule.FindStringSubmatch(command)
		dataType := "path"
		if len(parts) > 1 {
			if len(parts) > 2 {
				dataType = parts[1]
			}
			idx, _ := strconv.Atoi(parts[len(parts)-1])
			if idx < len(matrix.Matrix) {
				switch dataType {
				case "elem":
					renderOutput(idx, func(itemID int) []string {
						return []string{defaultArena.StringifyNode(itemID)}
					})
				case "path":
					renderOutput(idx, func(itemID int) []string {
						return []string{defaultArena.Chain(itemID, 0).XPath()}
					})
				case "html":
					renderOutput(idx, func(itemID int) []string {
						data, _ := defaultArena.RenderString(itemID)
						return []string{data}
					})
				case "text":
					renderOutput(idx, func(itemID int) []string {
						return defaultArena.StringifyInformation(itemID)
					})
				}
			}
		}
	}
}
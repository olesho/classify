package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/nathan-fiscaletti/consolesize-go"
	"github.com/olesho/classify/stream"
)

//var defaultArena *arena.Arena
var engine *stream.Engine
var matrix *stream.Matrix

func funcReset(command string) {
	matrix = nil
	engine = stream.NewEngine()
	//defaultArena = arena.NewArena()
}

func funcClear(command string) {
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
	for i, series := range matrix.Matrix[groupIdx].Nonuniform().Matrix {
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

func funcShow(command string) {
	if matrix == nil {
		//matrix = cluster.Extract(defaultArena)
		//matrix = cluster.ExtractOptimized(defaultArena)
		matrix = engine.Run(0, 4)
		fmt.Printf("Loaded succesfully. Total groups: %v\n", len(matrix.Matrix))
	}
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
						return []string{engine.Arena.StringifyNode(itemID)}
					})
				case "path":
					renderOutput(idx, func(itemID int) []string {
						return []string{engine.Arena.Chain(itemID, 0).XPath()}
					})
				case "html":
					renderOutput(idx, func(itemID int) []string {
						data, _ := engine.Arena.RenderString(itemID)
						return []string{data}
					})
				case "text":
					renderOutput(idx, func(itemID int) []string {
						return engine.Arena.StringifyInformation(itemID)
					})
				}
			}
		}
	}
}

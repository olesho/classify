package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/olesho/classify/arena"
	"github.com/olesho/classify/cluster"
	"os"
	"os/exec"
	"strconv"
	"runtime"
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

func renderOutput(groupIdx int, render func(itemID int) string) {
	w := prompt.NewStdoutWriter()
	rowsCnt := 0
	for _, series := range matrix.Matrix[groupIdx].Matrix {
		itemCnt := 0
		for _, item := range series {
			w.SetColor(prompt.DarkGray, prompt.White, true)
			w.WriteStr("\n"+render(item.Id)+"\n")
			w.Flush()

			w.SetColor(prompt.Red, prompt.DefaultColor, true)
			w.WriteStr("\n")
			w.Flush()
			//w.WriteStr("--------------------------------------------------------------------")
			itemCnt++
		}
		w.SetColor(prompt.Red, prompt.DefaultColor, true)
		w.WriteStr(fmt.Sprintf("\n        %v items total\n", itemCnt))
		w.WriteStr("====================================================================\n")
		w.Flush()
		rowsCnt++
	}
	fmt.Printf("        %v rows total\n", rowsCnt)
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
					renderOutput(idx, func(itemID int) string {
						return defaultArena.StringifyNode(itemID)
					})
				case "path":
					renderOutput(idx, func(itemID int) string {
						return defaultArena.Chain(itemID, 0).XPath()
					})
				case "html":
					renderOutput(idx, func(itemID int) string {
						data, _ := defaultArena.RenderString(itemID)
						return data
					})
				case "text":
					renderOutput(idx, func(itemID int) string {
						return defaultArena.StringifyInformation(itemID)
					})
				}
			}
		}
	}
}
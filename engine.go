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

var engine *stream.Storage
var matrix []stream.Series

func funcReset(command string) {
	matrix = nil
	engine = stream.NewStorage()
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
	bw.Open(prompt.Red, prompt.White, fmt.Sprintf("Total groups:%v", len(matrix[groupIdx].Matrix)))
	for i, series := range matrix[groupIdx].Nonuniform().Matrix {
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

func funcRun(command string) {
	//parts := runRule.FindStringSubmatch(command)
	//windowLength := 0
	//if len(parts) > 0 {
	//	windowLength, _ = strconv.Atoi(parts[1])
	//}
	//engine.SetWindowSize(windowLength)
	matrix = engine.RunAsync()
	fmt.Printf("succesfully loaded %v groups\n", len(matrix))
}

func funcShow(command string) {
	if matrix == nil {
		matrix = engine.Run()
		fmt.Printf("succesfully loaded %v groups\n" +
			"" +
			"", len(matrix))
	}
	if matrix != nil {
		parts := showRule.FindStringSubmatch(command)
		dataType := "path"
		if len(parts) > 1 {
			if len(parts) > 2 {
				dataType = parts[1]
			}
			idx, _ := strconv.Atoi(parts[len(parts)-1])
			if idx < len(matrix) {
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
				case "info":
					fmt.Println(engine.Arena.Chain(matrix[idx].Matrix[0][0].Id, 0).XPath())
					fmt.Printf("rows in group: %v\n", len(matrix[idx].Matrix))
					fmt.Printf("nodes in cluster: %v\n", len(matrix[idx].Group.Clusters))
					fmt.Printf("group volume: %v\n", matrix[idx].Group.GroupVolume)
					fmt.Printf("group size: %v\n", matrix[idx].Group.Size)
					fmt.Printf("volume: %v\n", matrix[idx].Group.Volume)
					fmt.Printf("wholesome volume: %v\n", matrix[idx].Group.WholesomeVolume)
				}
			}
		}
	}
}

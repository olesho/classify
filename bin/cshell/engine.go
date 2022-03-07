package main

import (
	"fmt"
	"github.com/olesho/classify/sequence"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/nathan-fiscaletti/consolesize-go"
)

var engine *sequence.RootCluster
var matrix []*sequence.Series

//var matrix []*stream.Series

func funcReset(command string) {
	matrix = nil
	engine = sequence.NewRootCluster()
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
	bw.Open(prompt.Red, prompt.White, fmt.Sprintf("Total groups:%v", len(matrix[groupIdx].TransposedNodes)))
	for i, transposedNodes := range matrix[groupIdx].TransposedNodes {
		bw.Open(prompt.Brown, prompt.White, fmt.Sprintf("Group %v", i+1))
		for nodeIndex, node := range transposedNodes {
			subItems := render(node.Id)
			bw.Open(prompt.Cyan, prompt.White, fmt.Sprintf("Item %v", nodeIndex+1))
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

func renderFields(groupIdx int) {
	cols, _ := consolesize.GetConsoleSize()
	bw := NewBlockWriter(cols+1, 1, 1)
	bw.Open(prompt.Red, prompt.White, fmt.Sprintf("Total groups:%v", len(matrix[groupIdx].TransposedValues)))
	for i, fields := range matrix[groupIdx].TransposedValues {
		bw.Open(prompt.Brown, prompt.White, fmt.Sprintf("Group %v", i+1))
		for fieldIndex, field := range fields {
			bw.Open(prompt.Cyan, prompt.White, fmt.Sprintf("Item %v", fieldIndex+1))
			bw.WriteText(prompt.Black, prompt.White, false, field)
			bw.Close()
		}
		bw.Close()
	}
	bw.Close()
	return
}

func renderPaths(groupIdx int) {
	cols, _ := consolesize.GetConsoleSize()
	bw := NewBlockWriter(cols+1, 1, 1)
	bw.Open(prompt.Red, prompt.White, fmt.Sprintf("Total groups:%v", len(matrix[groupIdx].TransposedValues)))

	for fieldIndex, xpath := range matrix[groupIdx].XPaths() {
		bw.Open(prompt.Cyan, prompt.White, fmt.Sprintf("Field %v", fieldIndex+1))
		bw.WriteText(prompt.Black, prompt.White, false, xpath)
		bw.Close()
	}
	bw.Close()
	return
}

func funcRun(command string) {
	parts := runRule.FindStringSubmatch(command)
	if len(parts) > 1 {
		numStr := strings.TrimSpace(parts[1])
		if windowLength, err := strconv.Atoi(numStr); err == nil && windowLength > 0 {
			engine.SetLimit(windowLength)
			fmt.Println("limit set to", windowLength)
		}
	}
	matrix = engine.Batch().Results()
	fmt.Printf("succesfully loaded %v groups\n", len(matrix))
}

func funcShow(command string) {
	if matrix == nil {
		matrix = engine.Batch().Results()
		fmt.Printf("succesfully loaded %v groups\n"+
			""+
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
					//renderOutput(idx, func(itemID int) []string {
					//	return []string{engine.Arena.Chain(itemID, 0).XPath()}
					//})
					renderPaths(idx)
				case "html":
					renderOutput(idx, func(itemID int) []string {
						data, _ := engine.Arena.RenderString(itemID)
						return []string{data}
					})
				case "text":
					renderOutput(idx, func(itemID int) []string {
						return engine.Arena.StringifyInformation(itemID)
					})
				case "fields":
					renderFields(idx)
				case "info":
					fmt.Println(engine.Arena.Chain(matrix[idx].TransposedNodes[0][0].Id, 0).XPath())
					fmt.Printf("rows in group: %v\n", len(matrix[idx].TransposedNodes))
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

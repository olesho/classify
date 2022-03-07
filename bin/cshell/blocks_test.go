package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/nathan-fiscaletti/consolesize-go"
	"testing"
)

func TestTest(t *testing.T) {
	cols, _ := consolesize.GetConsoleSize()
	fmt.Println(cols)
	bw := NewBlockWriter(cols+1, 1, 1)

	bw.Open(prompt.Red, prompt.White, "Total groups:14")

	bw.Open(prompt.Brown, prompt.Cyan, "Group 1")
	bw.Open(prompt.Blue, prompt.Yellow, "Item 1")
	bw.WriteText(prompt.Black, prompt.White, false, "How  to take the bus and train safely")
	bw.WriteText(prompt.Black, prompt.White, false, "Lesser-known  ways to reduce your risk of catching the coronavirus")
	bw.WriteText(prompt.Black, prompt.White, false, "Future")
	bw.Close()
	bw.Open(prompt.Blue, prompt.Yellow, "Item 2")
	bw.WriteText(prompt.Black, prompt.White, false, "How to take the bus and train safely")
	bw.Close()
	bw.Close()

	bw.Open(prompt.Brown, prompt.Cyan, "Group 2")
	bw.Open(prompt.Blue, prompt.Yellow, "Item 1")
	bw.WriteText(prompt.Black, prompt.White, false, "A  mystery in Asia's forgotten desert")
	bw.WriteText(prompt.Black, prompt.White, false, "For decades,  thousands of stone structures have baffled archaeologists")
	bw.WriteText(prompt.Black, prompt.White, false, "Travel")
	bw.Close()
	bw.Open(prompt.Blue, prompt.Yellow, "Item 2")
	bw.WriteText(prompt.Black, prompt.White, false, "A mystery in Asia's forgotten desert")
	bw.Close()
	bw.Close()

	bw.Open(prompt.Brown, prompt.Cyan, "Sub-title1")
	bw.Open(prompt.Blue, prompt.Yellow, "Sub-sub-title1")
	//bw.WriteText(prompt.Black, prompt.White, false, "Lorem ipsum dolor sit amet, consectetur adipiscing elit,\n sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\nUt enim ad minim veniam, quis nostrud exercitation ullamco laboris\n nisi ut aliquip ex ea commodo consequat.")
	//bw.WriteText(prompt.Black, prompt.White, false, "Lorem ipsum dolor sit amet, consectetur adipiscing elit,\n sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\nUt enim ad minim veniam, quis nostrud exercitation ullamco laboris\n nisi ut aliquip ex ea commodo consequat.")
	//bw.WriteText(prompt.Black, prompt.White, false, "Lorem ipsum dolor sit amet, consectetur adipiscing elit,\n sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\nUt enim ad minim veniam, quis nostrud exercitation ullamco laboris\n nisi ut aliquip ex ea commodo consequat.")
	bw.Close()
	bw.Close()

	bw.Close()
}

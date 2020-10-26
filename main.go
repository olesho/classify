package main

import (
	"fmt"
	"github.com/olesho/classify/sequence"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

var history []string

func main() {
	f, err := os.Open("history.log")
	if err != nil {
		log.Println(err)
	}
	bts, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	history = strings.Split(string(bts), "\n")
	engine = sequence.NewRootCluster()
	for {
		input := prompt.Input("> ", completer, prompt.OptionHistory(history))
		history = append(history, input)
		ok := false
		for _, r := range commands {
			if r.Regexp.MatchString(input) {
				r.Func(input)
				ok = true
				break
			}
		}
		if !ok {
			fmt.Printf("commant not found: %v\n", input)
		}
	}
}

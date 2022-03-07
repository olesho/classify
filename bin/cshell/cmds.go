package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
)

type Cmd struct {
	Regexp  *regexp.Regexp
	Func    func(string)
	Suggest prompt.Suggest
}

var qRule = regexp.MustCompile(`^q`)
var exitRule = regexp.MustCompile(`^exit`)
var fileRule = regexp.MustCompile(`^file\s+?(.+)`)
var helpRule = regexp.MustCompile(`\?`)
var runRule = regexp.MustCompile(`^run(\s+\d+)?`)
var showRule = regexp.MustCompile(`^show\s+(?:([a-z]+)?)\s+?([0-9]+)$`)
var clearRule = regexp.MustCompile(`^clear`)
var resetRule = regexp.MustCompile(`^reset`)
var rRule = regexp.MustCompile(`^r`)

var commands = []Cmd{
	{
		Regexp:  qRule,
		Func:    FuncExit,
		Suggest: prompt.Suggest{Text: "q", Description: "exit terminal"},
	},
	{
		Regexp:  exitRule,
		Func:    FuncExit,
		Suggest: prompt.Suggest{Text: "exit", Description: "exit terminal"},
	},
	{
		Regexp: fileRule,
		Func: func(command string) {
			parts := fileRule.FindStringSubmatch(command)
			if len(parts) > 1 {
				if engine != nil {
					err := engine.LoadFile(parts[1])
					if err != nil {
						log.Println(err)
						return
					}
					fmt.Printf("%v nodes total\n", len(engine.Arena.List))
				} else {
					log.Println("no context")
				}
			}
		},
		Suggest: prompt.Suggest{Text: "file ", Description: `file "file name" - load file`},
	},
	{
		Regexp:  webRule,
		Func:    funcWeb,
		Suggest: prompt.Suggest{Text: "web ", Description: `web "url" - load web page from URL`},
	},
	{
		Regexp:  wRule,
		Func:    funcWeb,
		Suggest: prompt.Suggest{Text: "w ", Description: `w "url" - load web page from URL`},
	},
	{
		Regexp:  chromeRule,
		Func:    funcChrome,
		Suggest: prompt.Suggest{Text: "chrome ", Description: `chrome "url" - load webpage from URL using Headless Chrome engine`},
	},
	{
		Regexp: helpRule,
		Func: func(command string) {
			fmt.Println(`
			file "file name" - load file
			http "url" - load from URL
			show [data type] [index] - shows an element with [index] from current list; possible data types: "path" (default), "content"  
			q/exit - exit terminal
		`)
		},
		Suggest: prompt.Suggest{Text: "help", Description: "help"},
	},
	{
		Regexp:  runRule,
		Func:    funcRun,
		Suggest: prompt.Suggest{Text: "run 200", Description: "run N - starts analysis with running window of size N"},
	},
	{
		Regexp:  showRule,
		Func:    funcShow,
		Suggest: prompt.Suggest{Text: "show elem", Description: "show elem [group_index] - shows content for element with [group_index] from current list"},
	},
	{
		Regexp:  showRule,
		Func:    funcShow,
		Suggest: prompt.Suggest{Text: "show path", Description: "show path [group_index] - shows path for element with [group_index] from current list"},
	},
	{
		Regexp:  showRule,
		Func:    funcShow,
		Suggest: prompt.Suggest{Text: "show html", Description: "show html [group_index] - shows content for element with [group_index] from current list"},
	},
	{
		Regexp:  showRule,
		Func:    funcShow,
		Suggest: prompt.Suggest{Text: "show text", Description: "show text [group_index] - shows content for element group_index [index] from current list"},
	},
	{
		Regexp:  showRule,
		Func:    funcShow,
		Suggest: prompt.Suggest{Text: "show fields", Description: "show fields [group_index] - shows content for element group_index [index] from current list"},
	},
	{
		Regexp:  clearRule,
		Func:    funcClear,
		Suggest: prompt.Suggest{Text: "clear", Description: "clear output"},
	},
	{
		Regexp:  resetRule,
		Func:    funcReset,
		Suggest: prompt.Suggest{Text: "reset", Description: "reset context"},
	},
	{
		Regexp:  rRule,
		Func:    funcReset,
		Suggest: prompt.Suggest{Text: "r", Description: "reset context"},
	},
}

func completer(d prompt.Document) []prompt.Suggest {
	var r []prompt.Suggest
	for _, c := range commands {
		if c.Regexp.MatchString(d.Text) || strings.HasPrefix(c.Suggest.Text, d.Text) || d.Text == "" {
			r = append(r, c.Suggest)
		}
	}
	return r
}

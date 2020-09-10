package main

import (
	"bufio"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/olesho/classify/arena"
	"github.com/olesho/classify/cluster"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Cmd struct {
	Regexp *regexp.Regexp
	Func func(string)
	Suggest prompt.Suggest
}

var qRule = regexp.MustCompile(`q`)
var exitRule = regexp.MustCompile(`exit`)
var fileRule = regexp.MustCompile(`file\s+?(.+)`)
var httpRule = regexp.MustCompile(`http\s+?(.+)`)
var helpRule = regexp.MustCompile(`\?`)
var showRule = regexp.MustCompile(`show\s+(?:([a-z]+)?)\s+?([0-9]+)$`)
var clearRule = regexp.MustCompile(`clear`)

var commands = []Cmd{
	{
		Regexp: qRule,
		Func: FuncExit,
		Suggest: prompt.Suggest{Text: "q", Description: "exit terminal"},
	},
	{
		Regexp: exitRule,
		Func: FuncExit,
		Suggest: prompt.Suggest{Text: "exit", Description: "exit terminal"},
	},
	{
		Regexp: fileRule,
		Func: func(command string) {
			parts := fileRule.FindStringSubmatch(command)
			if len(parts) > 1 {
				f, err := os.Open(parts[1])
				if err != nil {
					log.Println(err)
					return
				}
				defer f.Close()
				reader := bufio.NewReader(f)
				n, err := html.Parse(reader)
				defaultArena = arena.NewArena()
				defaultArena.Load(*n)
				matrix = cluster.Extract(defaultArena)

				fmt.Printf("Loaded succesfully. Total groups: %v\n", len(matrix.Matrix))
			}
		},
		Suggest: prompt.Suggest{Text: "file ", Description: `file "file name" - load file`},
	},
	{
		Regexp: httpRule,
		Func: func(command string) {
			parts := httpRule.FindStringSubmatch(command)
			if len(parts) > 1 {
				u := parts[1]
				if !(strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
					u = "https://" + u
				}
				r, err := http.Get(u)
				if err != nil {
					log.Println(err)
					return
				}
				reader := bufio.NewReader(r.Body)
				n, err := html.Parse(reader)
				defaultArena = arena.NewArena()
				defaultArena.Load(*n)
				matrix = cluster.Extract(defaultArena)

				fmt.Printf("Loaded succesfully. Total groups: %v\n", len(matrix.Matrix))
			}
		},
		Suggest: prompt.Suggest{Text: "http ", Description: `http "url" - load web page from URL`},
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
		Regexp: showRule,
		Func: FuncShow,
		Suggest: prompt.Suggest{Text: "show elem", Description: "show elem [group_index] - shows content for element with [group_index] from current list"},
	},
	{
		Regexp: showRule,
		Func: FuncShow,
		Suggest: prompt.Suggest{Text: "show path", Description: "show path [group_index] - shows path for element with [group_index] from current list"},
	},
	{
		Regexp: showRule,
		Func: FuncShow,
		Suggest: prompt.Suggest{Text: "show html", Description: "show html [group_index] - shows content for element with [group_index] from current list"},
	},
	{
		Regexp: showRule,
		Func: FuncShow,
		Suggest: prompt.Suggest{Text: "show text", Description: "show text [group_index] - shows content for element group_index [index] from current list"},
	},
	{
		Regexp: clearRule,
		Func: FuncClear,
		Suggest: prompt.Suggest{Text: "clear", Description: "clear output"},
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
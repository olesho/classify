package main

import (
	"encoding/json"
	"fmt"
	"github.com/olesho/classify/sequence"
	"os"
	"strconv"
)

type stats struct {
	GroupsCount int `json:"groups_count"`
	GroupFieldsCount int `json:"group_fields_count"`
}

type jsonResp struct {
	Fields [][]string `json:"fields"`
	Stats stats `json:"stats"`
}

var reversedOrder bool

func main() {
	var rank int
	for _, arg := range os.Args {
		if val, err := strconv.Atoi(arg); err == nil && val > 0 {
			rank = val
		} else if arg == "-r" {
			reversedOrder = true
		}
	}

	cluster := sequence.NewRootCluster()
	cluster.Load(os.Stdin)
	series := cluster.Batch().Results()
	if len(series) > rank {
		var i int
		for i = 0; i < len(os.Args); i++ {
			if os.Args[i] == "-xpath" {
				xpathOutput(series[rank])
				break
			} else if os.Args[i] == "-json" {
				jsonOutput(series[rank])
				break
			} else if os.Args[i] == "-text" {
				textOutput(series[rank])
				break
			}
		}

		if i == len(os.Args) {
			textOutput(series[rank])
		}
	}
}

func xpathOutput(s *sequence.Series) {
	for _, path := range s.XPaths() {
		fmt.Println(path)
	}
}

func jsonOutput(s *sequence.Series) {
	err := json.NewEncoder(os.Stdout).Encode(jsonResp{
		Fields: s.TransposedFields,
		Stats: stats{
			GroupsCount: len(s.TransposedFields),
			GroupFieldsCount: len(s.TransposedFields[0]),
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func textOutput(s *sequence.Series) {
	if reversedOrder {
		for i := len(s.TransposedFields)-1; i > -1; i-- {
			fields := s.TransposedFields[i]
			for _, field := range fields {
				fmt.Println(field)
			}
			fmt.Println("_________________________________________________")
		}
	} else {
		for _, fields := range s.TransposedFields {
			for _, field := range fields {
				fmt.Println(field)
			}
			fmt.Println("_________________________________________________")
		}
	}
}

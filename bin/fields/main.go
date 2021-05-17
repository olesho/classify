package main

import (
	"github.com/olesho/classify/sequence"
	"os"
	"strconv"
)

type stats struct {
	GroupsCount      int `json:"groups_count"`
	GroupFieldsCount int `json:"group_fields_count"`
}

type jsonResp struct {
	Fields [][]string `json:"fields"`
	Stats  stats      `json:"stats"`
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
				series[rank].ToXPath(os.Stdout)
				break
			} else if os.Args[i] == "-json" {
				series[rank].ToJSON(os.Stdout)
				break
			} else if os.Args[i] == "-text" {
				series[rank].ToText(os.Stdout)
				break
			} else if os.Args[i] == "-csv" {
				series[rank].ToCSV(os.Stdout)
				break
			}
		}

		if i == len(os.Args) {
			if reversedOrder {
				series[rank].ToTextReversed(os.Stdout)
				return
			}
			series[rank].ToText(os.Stdout)
		}
	}
}

package main

import (
	"fmt"
	"github.com/olesho/classify/sequence"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9876"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			typeString := r.URL.Query().Get("type")
			if typeString == "" {
				typeString = "json"
			}
			rankString := r.URL.Query().Get("rank")
			rank := 0
			if val, err := strconv.Atoi(rankString); err == nil && val > 0 {
				rank = val
			}

			err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
			if err != nil {
				http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
				return
			}

			cluster := sequence.NewRootCluster()
			// Media files

			for _, v := range r.MultipartForm.File {
				for _, h := range v {
					file, err := h.Open()
					if err != nil {
						http.Error(w, "failed to open file", http.StatusBadRequest)
						return
					}

					err = cluster.Load(file)
					if err != nil {
						http.Error(w, "failed to parse file", http.StatusBadRequest)
						return
					}

					err = file.Close()
					if err != nil {
						http.Error(w, "failed to close file", http.StatusBadRequest)
						return
					}
				}
			}

			series := cluster.Batch().Results()
			switch typeString {
			case "xpath":
				xpathOutput(series[rank], w)
			case "json":
				err = jsonOutput(series[rank], w)
				if err != nil {
					http.Error(w, "unable to write json", http.StatusBadRequest)
					return
				}
			case "text":
				textOutput(series[rank], w)
			case "csv":
				err = csvOutput(series[rank], w)
				if err != nil {
					http.Error(w, "unable to write csv", http.StatusBadRequest)
					return
				}
			default:
				textOutput(series[rank], w)
			}
			return
		}
	})

	fmt.Printf("listening on :%v\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		panic(err)
	}
}

func xpathOutput(s *sequence.Series, w http.ResponseWriter) {
	w.Header().Set("content-type", "text/plain")
	for _, path := range s.XPaths() {
		w.Write([]byte(fmt.Sprintln(path)))
	}
}

func jsonOutput(s *sequence.Series, w http.ResponseWriter) error {
	w.Header().Set("content-type", "application/json")
	return s.ToJSON(w)
}

func csvOutput(s *sequence.Series, w http.ResponseWriter) error {
	w.Header().Set("content-type", "text/csv")
	return s.ToCSV(w)
}

func textOutput(s *sequence.Series, w http.ResponseWriter) {
	w.Header().Set("content-type", "text/plain")
	s.ToText(w)
}

type stats struct {
	GroupsCount      int `json:"groups_count"`
	GroupFieldsCount int `json:"group_fields_count"`
}

type jsonResp struct {
	Fields [][]sequence.Field `json:"fields"`
	Stats  stats              `json:"stats"`
}

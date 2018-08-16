package main

import (
	"golang.org/x/net/html"
	"github.com/olesho/classify"
	htmlwords "github.com/olesho/go-words/html"
	"net/http"
	"github.com/olesho/go-words"
	"net/url"
	"encoding/json"
	"strconv"
	"log"
	"os"
)



func main() {
	http.Handle("/", http.StripPrefix("/static", http.FileServer(http.Dir("/static"))))

	http.HandleFunc("/parse/values", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		if r.Method == "GET" {
			u, err := url.Parse(r.RequestURI)
			if err != nil {
				json.NewEncoder(w).Encode(err)
				return
			}
			sourceUrl := u.Query().Get("url")
			rankStr := u.Query().Get("rank")
			var rank int
			if rankStr != "" {
				rank, err = strconv.Atoi(rankStr)
			}
			results, err := valuesFromWebsite(sourceUrl, rank)
			if err != nil {
				json.NewEncoder(w).Encode(err)
				return
			}
			err = json.NewEncoder(w).Encode(results)
			if err != nil {
				json.NewEncoder(w).Encode(err)
				return
			}
			return
		}
		if r.Method == "POST" {

		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal(err)
	}
}

/*
func valuesFromFiles(fileNames []string, rank int) ([][]string, error)  {
	cl := htmlwords.NewHtmlProcessor()

	for _, fName := range fileNames {
		f, err := os.Open(fName)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		n, err := html.Parse(f)
		if err != nil {
			return nil, err
		}

		arena := classify.NewArena(*n)
		for _, n := range arena.List {
			cl.Next(htmlwords.Node(n))
		}
	}

	bestWord := htmlwords.Word(cl.Vocabulary()[rank].Word)
	text := htmlwords.Word(cl.Text)
	return htmlwords.FilterValues(htmlwords.FindValues(bestWord, text)), nil
}
*/

func valuesFromWebsite(u string, rank int) ([][]string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
		n, err := html.Parse(resp.Body)
		if err != nil {
			return nil, err
		}

		arena := classify.NewArena(*n)
		cl := htmlwords.NewHtmlProcessor()
		for _, n := range arena.List {
			cl.Next(htmlwords.Node(n))
		}

		cl.SortVocabulary()

		bestWord := htmlwords.Word(cl.Vocabulary()[rank].Word)
		text := htmlwords.Word(cl.Text)
		return htmlwords.FilterValues(htmlwords.FindValues(bestWord, text)), nil
	}
	return nil, nil
}

func Tags(nn []classify.Node) []string {
	res := make([]string, len(nn))
	for i, n := range nn {
		res[i] = Tag(n)
	}
	return res
}

func Tag(n classify.Node) string {
	if n.Type == html.ElementNode {
		return n.Data
	} else if n.Type == html.TextNode {
		return "TEXT"
	} else {
		return "OTHER"
	}
}

func replace(n classify.Node) words.Str {
	return words.Str(Tag(n))
}

func toString(w []words.Comparable) string {
	s := ""
	for _, c := range w {
		s += " " + string(c.(words.Str))
	}
	return s
}
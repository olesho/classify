package main

import (
	"golang.org/x/net/html"
	"github.com/olesho/classify"
	htmlwords "github.com/olesho/classify/go-words/html"
	"net/http"
	"github.com/olesho/classify/go-words"
	"net/url"
	"encoding/json"
	"strconv"
	"log"
	"os"
	"fmt"
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

	http.HandleFunc("/parse/fields", func(w http.ResponseWriter, r *http.Request) {
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
			results, err := fieldsFromWebsite(sourceUrl, rank)
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
	_, cl, err := makeProcessor(u)
	if err != nil {
		return nil, err
	}
	bestWord := htmlwords.Word(cl.Vocabulary()[rank].Word)
	text := htmlwords.Word(cl.Text)
	return htmlwords.FilterValues(htmlwords.FindValues(bestWord, text)), nil
}

func fieldsFromWebsite(u string, rank int) (words.Word, error) {
	arena, cl, err := makeProcessor(u)
	if err != nil {
		return nil, err
	}

	var prev classify.Node
	bestWord := cl.Vocabulary()[rank].Word
	text := htmlwords.Word(cl.Text)
	positions := words.FindPositions(words.Word(bestWord), words.Word(text))
	for _, p := range positions {
		currentWord := text[p: p+len(bestWord)]
		for i, n := range currentWord {
			nw := classify.Node(n.(htmlwords.Node))
			if i > 0 {
				//prev := classify.Node(words.Word(currentWord)[i-1].(htmlwords.Node))
				if !arena.HasDescendant(prev.Id, nw.Id) {
					prev = nw
					render, _ := arena.RenderString(nw.Id)
					fmt.Println(render)
					fmt.Println("-------------------------------------------------------")
				} else {
					//fmt.Println("|_")
				}
			} else {
				prev = nw
				render, _ := arena.RenderString(nw.Id)
				fmt.Println(render)
				fmt.Println("-------------------------------------------------------")
			}
		}
		fmt.Println("=========================================================")
	}

	/*
	for _, n := range words.Word(bestWord) {
		if i > 0 {
			prev := classify.Node(words.Word(bestWord)[i-1].(htmlwords.Node))
			if !arena.HasParent(nw.Id, prev.Id) {
				fmt.Println(arena.RenderString(nw.Id))
			}
		} else {
			fmt.Println(arena.RenderString(nw.Id))
		}
	}
	*/
	//nw := classify.Node(n.(htmlwords.Node))
	//fmt.Println(arena.RenderString(nw.Id))
	return bestWord, nil
}

func makeProcessor(u string) (*classify.Arena , *htmlwords.HtmlProcessor, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, nil, err
	}
	if resp != nil {
		defer resp.Body.Close()
		n, err := html.Parse(resp.Body)
		if err != nil {
			return nil, nil, err
		}

		arena := classify.NewArena(*n)
		cl := htmlwords.NewHtmlProcessor()
		for _, n := range arena.List {
			cl.Next(htmlwords.Node(n))
		}

		cl.SortVocabulary()
		return arena, cl, nil
	}
	return nil, nil, nil
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
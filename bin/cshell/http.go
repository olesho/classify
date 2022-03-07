package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/net/html/charset"

	"golang.org/x/text/encoding/htmlindex"
)

var webRule = regexp.MustCompile(`^web\s+?(.+)`)
var wRule = regexp.MustCompile(`^w\s+?(.+)`)
var chromeRule = regexp.MustCompile(`^chrome\s+?(.+)`)

var charsetRule = regexp.MustCompile(`charset=([^()<>@,;:\"/[\]?.=\s]*)`)

func detectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

func decode(body io.Reader, charset string) (io.Reader, error) {
	if charset == "" {
		charset = detectContentCharset(body)
	}
	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}

	if name, _ := htmlindex.Name(e); name != "utf-8" {
		body = e.NewDecoder().Reader(body)
	}

	return body, nil
}

func getContentCharset(r *http.Response) string {
	contentType := r.Header.Get("content-type")
	s := charsetRule.FindAllStringSubmatch(contentType, 1)
	if len(s) > 0 {

		if len(s[0]) > 1 {
			return s[0][1]
		}
	}
	return ""
}

func funcWeb(command string) {
	parts := webRule.FindStringSubmatch(command)
	if len(parts) > 1 {
		u := parts[1]
		if !(strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")) {
			u = "https://" + u
		}

		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			log.Println(err)
			return
		}

		r, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}

		reader, err := decode(r.Body, getContentCharset(r))
		if err != nil {
			log.Println(err)
			return
		}
		if engine != nil {
			err := engine.Load(reader)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("%v nodes total\n", len(engine.Arena.List))
		} else {
			log.Println("empty context")
		}
	}
}

func funcChrome(command string) {
	parts := chromeRule.FindStringSubmatch(command)
	if len(parts) > 1 {
		u := parts[1]
		cmd := exec.Command(`google-chrome`, `--headless`, `--dump-dom`, u)
		data, err := cmd.Output()
		if err != nil {
			log.Println(err)
		}

		if engine != nil {
			err := engine.Load(bytes.NewBuffer(data))
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("%v nodes total\n", len(engine.Arena.List))
		} else {
			log.Println("empty context")
		}
	}
}

package main

import (
	"log"
	"os"
	"strings"
)

func FuncExit(command string) {
	f, err := os.Create("history.log")
	if err != nil {
		log.Println(err)
	}
	_, err = f.WriteString(strings.Join(history, "\n"))
	if err != nil {
		log.Println(err)
	}
	err = f.Close()
	if err != nil {
		log.Println(err)
	}
	os.Exit(0)
}

package main

import (
	"fmt"
	"os"
	"strings"
)

var exit_code chan int

func usage() {
	fmt.Println("Usage: yes [STRING]...")
	fmt.Println("  or:  yes OPTION")
	exit_code <- 1
}

func loop(msg string) {
	for {
	    _, err := fmt.Println(msg)
		if err != nil {
			exit_code <- 130
			break;
		}
	}
}

func main() {
	exit_code = make (chan int)
	msg := "y"
	run := true
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help":
			go usage()
			run = false;
		default:
			msg = strings.Join(os.Args[1:], " ")
		}
	}
	if run {
		loop(msg)
	}
	os.Exit(<-exit_code)
}

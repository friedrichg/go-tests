package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
        "bufio"
)

func usage(exit_code *chan int) {
	fmt.Println("Usage: yes [STRING]...")
	fmt.Println("  or:  yes OPTION")
	*exit_code <- 1
}

func loop(cont *bool, msg string, exit_code *chan int) {
        o := bufio.NewWriter(os.Stdout)
	for *cont {
		_, err := fmt.Fprintln(o,msg)
		if err != nil {
			*exit_code <- 120
			break
		}
	}
	o.Flush()
}

func check_signals(exit_code *chan int, sigs *chan os.Signal) {
	<-*sigs
	*exit_code <- 130
}

func main() {
	msg := "y"
	run := true
	exit_code := make(chan int)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go check_signals(&exit_code, &sigs)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help":
			go usage(&exit_code)
			run = false
		default:
			msg = strings.Join(os.Args[1:], " ")
		}
	}
	go loop(&run, msg, &exit_code)
	os.Exit(<-exit_code)
}

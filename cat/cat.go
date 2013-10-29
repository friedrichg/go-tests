package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	fmt.Println("Usage: cat [OPTION]... [FILE]...")
	//TODO: support -n
	os.Exit(0)
}

func cat(fi *os.File) {
	out := fmt.Print
	buf := make([]byte, 1024)
	for {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		out(string(buf[:n]))
	}

}

func main() {
	var exit_code int
	l := fmt.Println
	args := os.Args[1:]
	if len(args) == 0 {
		args = make([]string, 1)
		args[0] = "-"
	}

	add := true
	//Add default Stdin if only args are passed
	for _, a := range args {
		if a == "-" || !strings.HasPrefix(a, "-") {
			add = false
			break
		}
	}
	if add {
		args = append(args, "-")
	}

	//TODO: Parse more flags (-n)
	for _, a := range args {
		switch {
		case a == "--help":
			usage()
		case a == "-":
			cat(os.Stdin)
		case strings.HasPrefix(a, "-"):
			usage()
		default:
			fi, err := os.Open(a)
			if err != nil {
				l("cat: " + a + ": No such file or directory")
				exit_code = 1
				continue
			}
			defer func() {
				if err := fi.Close(); err != nil {
					panic(err)
				}
			}()
			cat(fi)
		}
	}
	os.Exit(exit_code)
}

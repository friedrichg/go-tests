package main

import (
	"fmt"
	"io"
	"os"
)

func usage() {
	fmt.Println("Usage: cat [OPTION]... [FILE]...")
	//TODO: show flags supported
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
	filename := os.Args[1:]
	if len(filename) == 0 {
		filename = make([]string, 1)
		filename[0] = "-"
	}

	//TODO: Parse more flags
	for _, f := range filename {
		switch f {
		case "--help":
			usage()
		case "-":
			cat(os.Stdin)
		default:
			fi, err := os.Open(f)
			if err != nil {
				l("cat: " + f + ": No such file or directory")
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

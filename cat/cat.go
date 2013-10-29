package main

import (
	"fmt"
	"io"
	"os"
)

func usage() {
	fmt.Println("Usage: cat [OPTION]... [FILE]...")
	//TODO: show flags supported
	os.Exit(1)
}

func cat(file []string) {
	l := fmt.Println
	out := fmt.Print
	for _, f := range file {
		fi, err := os.Open(f)
		if err != nil {
			l("cat: " + f + ": No such file or directory")
			continue
		}
		defer func() {
			if err := fi.Close(); err != nil {
				panic(err)
			}
		}()
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
}

func main() {
	if len(os.Args) < 2 {
		usage()
	} else {
		//TODO: support more flags
		cat(os.Args[1:])
	}
}

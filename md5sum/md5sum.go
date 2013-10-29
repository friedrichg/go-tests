package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	l := fmt.Println
	l("Usage: md5sum [OPTION]... [FILE]...")
	l()
	os.Exit(0)
}

func encodeFile(fi *os.File) {
	reader := bufio.NewReader(fi)
	buf := make([]byte, 1024)
	h := md5.New()
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		_, err = h.Write(buf[:n])
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("%x", h.Sum(nil))
}

func closeFile(fi *os.File) {
	if err := fi.Close(); err != nil {
		panic(err)
	}
}

func main() {
	var exit_code int
	args := os.Args[1:]

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

	//TODO: Parse more flags
	l := fmt.Println
	for _, a := range args {
		switch {
		case a == "--help":
			usage()
		case a == "-":
			encodeFile(os.Stdin)
		case strings.HasPrefix(a, "-"):
			usage()
		default:
			fi, err := os.Open(a)
			if err != nil {
				l("md5sum: " + a + ": No such file or directory")
				exit_code = 1
				continue
			}
			defer closeFile(fi)
			encodeFile(fi)
		}
	}
	os.Exit(exit_code)
}

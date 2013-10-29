package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func usage() {
	l := fmt.Println
	l("Usage: base64 [OPTION]... [FILE]...")
	l("Base64 encode or decode FILE, or standard input, to standard output.")
	l()
	l("  -d, --decode          decode data")
	l("  -i, --ignore-garbage  when decoding, ignore non-alphabet characters")
	l()
	os.Exit(0)
}

func decodeFile(fi *os.File, ignore bool) {
	out := fmt.Print
	reader := bufio.NewReader(fi)
	decoder := base64.NewDecoder(base64.StdEncoding, reader)
	buf := make([]byte, 1023)
	for {
		n, err := decoder.Read(buf)
		if !ignore && err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		out(string(buf[:n]))
		if err != nil {
			panic(err)
		}
	}
}

func encodeFile(fi *os.File, ignore bool) {
	reader := bufio.NewReader(fi)
	encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	buf := make([]byte, 1023)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		encoder.Write(buf[:n])
		if err != nil {
			panic(err)
		}
	}
	encoder.Close()
	fmt.Println()
}

func closeFile(fi *os.File) {
	if err := fi.Close(); err != nil {
		panic(err)
	}
}

func main() {
	var exit_code int
	ignore := false
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
	transform := encodeFile
	for _, a := range args {
		switch {
		case a == "-i", a == "--ignore-garbage":
			ignore = true
		case a == "-d", a == "--decode":
			transform = decodeFile
		case a == "--help":
			usage()
		case a == "-":
			transform(os.Stdin, ignore)
		case strings.HasPrefix(a, "-"):
			usage()
		default:
			fi, err := os.Open(a)
			if err != nil {
				l("base64: " + a + ": No such file or directory")
				exit_code = 1
				continue
			}
			defer closeFile(fi)
			transform(fi, ignore)
		}
	}
	os.Exit(exit_code)
}

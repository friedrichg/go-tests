package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func usage() {
	fmt.Println("Usage: cat [OPTION]... [FILE]...")
	//TODO: support -n
	os.Exit(0)
}

func openURL(url string) int {
	l := fmt.Println
	resp, err := http.Get(url)
	if err != nil {
		l("Error getting webpage", err)
		return 1
	}
	switch resp.StatusCode {
	case 200:
		//l("200 OK")
	default:
		l("ERROR", resp.StatusCode)
		return 1
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l("Error reading body", err)
		return 1
	}
	randomFile, err := ioutil.TempFile("", "base64")
	if err != nil {
		l("Error creating temp file")
		return 1
	}
	defer randomFile.Close()
	defer os.Remove(randomFile.Name())
	//TODO: this does not support big chunks. Fix it
	_, err = randomFile.Write(body)
	if err != nil {
		l("Error writing to temp file")
		return 1
	}
	_, err = randomFile.Seek(0, os.SEEK_SET)
	if err != nil {
		l("Error restarting temp file")
		return 1
	}
	cat(randomFile)
	return 0
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
		case strings.HasPrefix(a, "http://"):
			openURL(a)
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

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func usage() {
	fmt.Println("Usage: head [OPTION]... [FILE]...")
	//TODO: support -n
	os.Exit(0)
}

func openURL(url string, maxlines int) int {
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
	head(randomFile, maxlines)
	return 0
}

func head(fi *os.File, maxlines int) {
	out := fmt.Print
	buf := make([]byte, 1024)
	var endofline = []byte("\n")
	for maxlines > 0 {
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		//Counts number of lines in buffer
		el := bytes.Count(buf[:n], endofline)

		//If number of lines surpasses the maximum we cut the rest
		if el > maxlines {
			begin := 0
			a := 0
			for i := 0; i < maxlines; i++ {
				fmt.Printf("begin=%d i=%d a=%d len=%d\n", begin, i, a, len(buf[begin:n]))
				t := buf[begin+1 : n]
				begin = bytes.IndexByte(t, 'a')
			}
			n = begin
			el = maxlines
		}
		out(string(buf[:n]))
		maxlines -= el
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

	regexpdashn := regexp.MustCompile("^-n[0-9]+$")
	maxlines := 10

	//TODO: Parse more flags (-n)
	for _, a := range args {
		switch {
		case a == "--help":
			usage()
		case a == "-":
			head(os.Stdin, maxlines)
		case regexpdashn.MatchString(a):
			maxlines, _ = strconv.Atoi(a[2:])
		case strings.HasPrefix(a, "-"):
			usage()
		case strings.HasPrefix(a, "http://"):
			if openURL(a, maxlines) != 0 {
				exit_code = 1
			}
		default:
			fi, err := os.Open(a)
			if err != nil {
				l("head: " + a + ": No such file or directory")
				exit_code = 1
				continue
			}
			defer func() {
				if err := fi.Close(); err != nil {
					panic(err)
				}
			}()
			head(fi, maxlines)
		}
	}
	os.Exit(exit_code)
}

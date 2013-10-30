package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func usage() {
	l := fmt.Println
	l("Usage: md5sum [OPTION]... [FILE]...")
	l()
	os.Exit(0)
}

func encodeFile(fi *os.File, name string) {
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
	fmt.Printf("%x %s\n", h.Sum(nil), name)
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
	encodeFile(randomFile, url)
	return 0
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
			encodeFile(os.Stdin, a)
		case strings.HasPrefix(a, "-"):
			usage()
		case strings.HasPrefix(a, "http://"):
			openURL(a)
		default:
			fi, err := os.Open(a)
			if err != nil {
				l("md5sum: " + a + ": No such file or directory")
				exit_code = 1
				continue
			}
			defer closeFile(fi)
			encodeFile(fi, a)
		}
	}
	os.Exit(exit_code)
}

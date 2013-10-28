package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func usage() {
	fmt.Println("Usage: wget [OPTION]... [URL]...")
	os.Exit(1)
}

func get(url string) {
	l := fmt.Println

	resp, err := http.Get(url)
	if err != nil {
		l("Error getting webpage", err)
		return
	}
	switch resp.StatusCode {
	case 200:
		l("200 OK")
	default:
		l("ERROR", resp.StatusCode)
		return
	}
	path := resp.Request.URL.Path
	i := strings.LastIndex(path, "/")
	path = path[i+1:]
	if path == "" {
		path = "index.html"
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l("Error reading body", err)
		return
	}
	err = ioutil.WriteFile(path, body, 0644)
	if err != nil {
		l("Error writing to file", err)
		return
	}
	l(path + " Saved")

}

func main() {
	if len(os.Args) < 2 {
		usage()
	} else {
		u := os.Args[1]
		uparsed, err := url.Parse(u)
		if err != nil {
			log.Println("La url es invalida", u)
			usage()
		} else {
			if uparsed.Scheme == "" {
				uparsed.Scheme = "http"
			}
			get(uparsed.String())
		}
	}

}

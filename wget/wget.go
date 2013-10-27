package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
)

func usage(){
	fmt.Println("Usage: wget [OPTION]... [URL]...")
	os.Exit(1)
}

func get(url string) {
	l:=fmt.Println
	resp, err:=http.Get(url)
	if err != nil {
		l("Error getting webpage",err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l("Error reading body")
		return
	}
	fmt.Println(string(body))
	
}

func main() {
	if (len(os.Args) < 2){
		usage()
	} else {
		url := os.Args[1]
		get(url)
	}
		
}

package main

import (
        "net/http"
        "log"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Header().Set("Connection", "keep-alive")
        w.Write([]byte("hello, world!\n"))
}
func main() {
        http.HandleFunc("/", HelloServer)
        log.Println("Serving at http://127.0.0.1:8080/")
        http.ListenAndServe(":8080", nil)
}

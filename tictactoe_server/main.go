package main

import (
    "github.com/aidarbiktimirov/tictactoe/api"

    "io/ioutil"
    "net/http"
    "fmt"
    "os"
)

func createStaticHandler(filename, mimetype string) func (w http.ResponseWriter, r *http.Request) {
    return func (w http.ResponseWriter, r *http.Request) {
    data, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", os.Args[1], filename))
        w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", mimetype))
        fmt.Fprintf(w, "%s", data)
    }
}

func main() {
    srv := api.NewServer()
    srv.HandleFunc("/out.js", createStaticHandler("out.js", "application/javascript"))
    srv.HandleFunc("/", createStaticHandler("index.html", "text/html"))
    srv.Run(":10000")
}

package main

import (
    "github.com/aidarbiktimirov/tictactoe/api"
)

func main() {
    srv := api.NewServer()
    // srv.HandleFunc("/tictactoe/", handleWebPage)
    // srv.HandleFunc("/static/", handleStaticFiles)
    srv.Run(":10000")
}

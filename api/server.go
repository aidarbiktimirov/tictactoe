package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
)

type Server struct {
    api *Api
    mux *http.ServeMux
}

func NewServer() *Server {
    result := &Server{}
    result.api = NewApi()
    result.mux = http.NewServeMux()
    result.mux.HandleFunc("/api/games/list/", result.handleListGames)
    result.mux.HandleFunc("/api/games/new/", result.handleNewGame)
    result.mux.HandleFunc("/api/games/show/", result.handleGetGameField)
    result.mux.HandleFunc("/api/games/update/", result.handleUpdateGame)
    return result
}

func (this *Server) HandleFunc(path string, handler func(w http.ResponseWriter, r *http.Request)) {
    this.mux.HandleFunc(path, handler)
}

func (this *Server) Run(addr string) error {
    server := &http.Server{Addr: addr, Handler: this.mux}
    return server.ListenAndServe()
}

func (this *Server) handleListGames(w http.ResponseWriter, r *http.Request) {
    res, _ := json.Marshal(this.api.ListGames())
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    fmt.Fprintf(w, "%s", res)
}

func (this *Server) handleNewGame(w http.ResponseWriter, r *http.Request) {
    id, err := this.api.NewGame(r.FormValue("playerX"), r.FormValue("playerO"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    } else {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        fmt.Fprintf(w, "%d", id)
    }
}

func (this *Server) handleGetGameField(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.FormValue("id"))
    x, _ := strconv.Atoi(r.FormValue("x"))
    y, _ := strconv.Atoi(r.FormValue("y"))
    width, _ := strconv.Atoi(r.FormValue("width"))
    height, _ := strconv.Atoi(r.FormValue("height"))
    fieldPart, err := this.api.GetGameFieldPart(id, x, y, width, height)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    } else {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        res, _ := json.Marshal(fieldPart)
        fmt.Fprintf(w, "%s", res)
    }
}

func (this *Server) handleUpdateGame(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.FormValue("id"))
    x, _ := strconv.Atoi(r.FormValue("x"))
    y, _ := strconv.Atoi(r.FormValue("y"))
    user := r.FormValue("username")
    state, err := this.api.UpdateGame(id, x, y, user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    } else {
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        res, _ := json.Marshal(state)
        fmt.Fprintf(w, "%s", res)
    }
}

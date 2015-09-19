package api

import (
    "errors"
    "sync"

    "github.com/aidarbiktimirov/tictactoe/engine"
)

type gameInfo struct {
    gameId int
    playerX, playerO string
    field *engine.InfiniteField
}

type Api struct {
    games []gameInfo
    gamesMutex *sync.Mutex
}

func NewApi() *Api {
    result := &Api{}
    result.gamesMutex = &sync.Mutex{}
    return result
}

func (this *Api) ListGames() []map[string]interface{} {
    result := make([]map[string]interface{}, len(this.games))
    for i := 0; i < len(result); i++ {
        result[i] = make(map[string]interface{})
        users := make([]interface{}, 2)
        users[0] = this.games[i].playerX
        users[1] = this.games[i].playerO
        result[i]["id"] = i
        result[i]["users"] = users
        switch this.games[i].field.GetState() {
        case engine.PlayerXVictory:
            result[i]["winner"] = "x"
        case engine.PlayerOVictory:
            result[i]["winner"] = "o"
        case engine.Draw:
            result[i]["winner"] = nil
        }
    }
    return result
}

func (this *Api) NewGame(playerX, playerO string) (int, error) {
    if playerX == playerO {
        return -1, errors.New("Players must be different")
    }
    this.gamesMutex.Lock()
    this.games = append(this.games, gameInfo{})
    id := len(this.games) - 1
    this.games[id].gameId = id
    this.games[id].playerX = playerX
    this.games[id].playerO = playerO
    this.games[id].field = engine.NewInfiniteField()
    this.gamesMutex.Unlock()
    return id, nil
}

func (this *Api) GetGameFieldPart(id int, x, y int, width, height int) ([][]interface{}, error) {
    if id < 0 || id >= len(this.games) {
        return nil, errors.New("Incorrect game id")
    }
    game := this.games[id].field
    if game == nil {
        return nil, errors.New("Incorrect game id")
    }

    result := make([][]interface{}, width)
    this.gamesMutex.Lock()
    for i := 0; i < width; i++ {
        result[i] = make([]interface{}, height)
        for j := 0; j < height; j++ {
            switch game.Get(x + i, y + j) {
            case engine.PlayerX:
                result[i][j] = "x"
            case engine.PlayerO:
                result[i][j] = "o"
            default:
                result[i][j] = nil
            }
        }
    }
    this.gamesMutex.Unlock()
    return result, nil
}

func (this *Api) UpdateGame(id int, x, y int, player string) (map[string]string, error) {
    if id < 0 || id >= len(this.games) {
        return nil, errors.New("Incorrect game id")
    }

    parseState := func(state engine.State, game gameInfo) map[string]string {
        res := make(map[string]string)
        switch state {
        case engine.PlayerXVictory:
            res["state"] = "victory"
            res["winner"] = game.playerX
        case engine.PlayerOVictory:
            res["state"] = "victory"
            res["winner"] = game.playerO
        case engine.Draw:
            res["state"] = "draw"
        case engine.Unknown:
            res["state"] = "in progress"
        }
        return res
    }

    game := this.games[id]
    switch player {
    case game.playerX:
        state, err := game.field.Set(x, y, engine.PlayerX)
        return parseState(state, game), err
    case game.playerO:
        state, err := game.field.Set(x, y, engine.PlayerO)
        return parseState(state, game), err
    default:
        return nil, errors.New("Incorrect player")
    }
}

package engine

import (
    "errors"
    "sync"
)

type Player int

const (
    NoPlayer Player = iota
    PlayerX
    PlayerO
)

type State int

const (
    Unknown State = iota
    Draw
    PlayerXVictory
    PlayerOVictory
)

type Fielder interface {
    Get(x, y int) Player
    Set(x, y int, cell Player) (State, error)
    LastTurn() Player
    GetState() State
}

// Fielder helpers

func validateTurn(field Fielder, x, y int, cell Player) error {
    if cell == NoPlayer {
        return errors.New("Can't set empty value")
    }
    if field.GetState() != Unknown {
        return errors.New("The game is already finished")
    }
    if cell == field.LastTurn() {
        return errors.New("Same player can't make turns twice in a row")
    }
    if field.Get(x, y) != NoPlayer {
        return errors.New("Can't fill non-empty cell")
    }
    return nil
}

func getWinnerState(player Player) State {
    switch player {
    case PlayerX:
        return PlayerXVictory
    case PlayerO:
        return PlayerOVictory
    default:
        return Unknown
    }
}

func getState(field Fielder, x, y int, width int) State {
    if winner := field.Get(x, y); winner != NoPlayer {
        for _, dir := range [...][2]int{{1, 0}, {0, 1}, {1, 1}} {
            sum := 1
            for _, m := range [...]int{-1, 1} {
                for i := 1; i < width && field.Get(x + i * dir[0] * m, y + i * dir[1] * m) == winner; i++ {
                    sum++
                }
            }
            if sum >= width {
                return getWinnerState(winner)
            }
        }
    }
    return Unknown
}

// Fixed 3x3 field

type Field3x3 struct {
    cells [3][3]Player
    lastTurn Player
    state State
    nonNoPlayerCells uint
    mutex *sync.Mutex
}

func NewField3x3() *Field3x3 {
    result := &Field3x3{}
    result.mutex = &sync.Mutex{}
    return result
}

func (this *Field3x3) Get(x, y int) Player {
    if x < 0 || x >= 3 || y < 0 || y >= 3 {
        return NoPlayer
    } else {
        return this.cells[x][y]
    }
}

func (this *Field3x3) LastTurn() Player {
    return this.lastTurn
}

func (this *Field3x3) GetState() State {
    return this.state
}

func (this *Field3x3) Set(x, y int, cell Player) (State, error) {
    this.mutex.Lock()
    if err := validateTurn(this, x, y, cell); err != nil {
        result := this.GetState()
        this.mutex.Unlock()
        return result, err
    }
    this.cells[x][y] = cell
    this.lastTurn = cell
    this.nonNoPlayerCells++
    this.state = getState(this, x, y, 3)
    if this.state == Unknown && this.nonNoPlayerCells == 9 {
        this.state = Draw
    }
    result := this.GetState()
    this.mutex.Unlock()
    return result, nil
}

// Infinite field

type InfiniteField struct {
    cells map[int]map[int]Player
    lastTurn Player
    state State
    mutex *sync.Mutex
}

func NewInfiniteField() *InfiniteField {
    result := &InfiniteField{}
    result.mutex = &sync.Mutex{}
    return result
}

func (this *InfiniteField) Get(x, y int) Player {
    if row, ok := this.cells[x][y]; ok {
        return row
    }
    return NoPlayer
}

func (this *InfiniteField) Set(x, y int, cell Player) (State, error) {
    this.mutex.Lock()
    if err := validateTurn(this, x, y, cell); err != nil {
        this.mutex.Unlock()
        result := this.GetState()
        return result, err
    }
    if this.cells == nil {
        this.cells = make(map[int]map[int]Player)
    }
    if _, ok := this.cells[x]; !ok {
        this.cells[x] = make(map[int]Player)
    }
    this.cells[x][y] = cell
    this.lastTurn = cell
    this.state = getState(this, x, y, 5)
    result := this.GetState()
    this.mutex.Unlock()
    return result, nil
}

func (this *InfiniteField) LastTurn() Player {
    return this.lastTurn
}

func (this *InfiniteField) GetState() State {
    return this.state
}

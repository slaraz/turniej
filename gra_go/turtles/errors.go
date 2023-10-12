package turtles

import (
	"fmt"
)

var ErrGameNotFound = fmt.Errorf("game not found")
var ErrBoardNotFound = fmt.Errorf("board not found")
var ErrInvalidMove = fmt.Errorf("invalid move")
var ErrInvalidPlayerNumber = fmt.Errorf("invalid player number")
var ErrInvalidCard = fmt.Errorf("invalid card")
var ErrInvalidColor = fmt.Errorf("invalid color")
var ErrPickTheColor = fmt.Errorf("pick the color")
var ErrGameIsOver = fmt.Errorf("game is over")
var ErrNoPlayers = fmt.Errorf("no players")

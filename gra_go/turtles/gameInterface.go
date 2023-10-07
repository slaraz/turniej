package turtles

import (
	"encoding/json"
	"strings"
)

type Move struct {
	CardSymbol string `json:"cardSymbol"`
	Color      string `json:"color"`
}
type GameStatus struct {
	Board       []Field `json:"board"`
	Cards       []Card  `json:"cards"`
	Winer       int     `json:"winer"`
	IsEnd       bool    `json:"isEnd"`
	TurtleColor string  `json:"turtleColor"`
}

func (game *Game) GetGameStatus(playerNumber int) (string, error) {
	if playerNumber-1 > len(game.players) || playerNumber-1 < 0 {
		return "", ErrInvalidPlayerNumber
	}
	status := GameStatus{
		Board:       game.board,
		Cards:       game.players[playerNumber-1].Cards,
		TurtleColor: string(game.players[playerNumber-1].Color),
		Winer:       game.winer, //IF WINER IS -1 THEN NO WINER
		IsEnd:       game.isEnd,
	}
	json, _ := json.Marshal(status)
	return string(json), nil
}
func (game *Game) Move(moveStr string, playerNumber int) error {
	playerNumber = playerNumber - 1
	if playerNumber > len(game.players) || playerNumber < 0 {
		return ErrInvalidPlayerNumber
	}
	move := Move{}
	err := json.Unmarshal([]byte(moveStr), &move)
	if err != nil {
		return err
	}
	move.CardSymbol = strings.ToUpper(move.CardSymbol)
	card, err := findCard(Symbol(move.CardSymbol))
	if err != nil {
		return err
	}
	move.Color = strings.ToLower(move.Color)
	color := getColor(move.Color)
	err = game.playCard(card, color, playerNumber)
	if err != nil {
		return err
	}
	if game.isEnd {
		return nil
	}
	return nil
}
func CreateNewGame(numberOfPlayers int) Game {
	game := Game{
		board:      CreateGameBoard("a"),
		deck:       CreateGameDeck("a"),
		players:    generatePlayers(numberOfPlayers),
		round:      0,
		playerTurn: 0,
	}
	game.dealTheCards()
	return game
}

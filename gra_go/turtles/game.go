package turtles

import (
	"log"
)

var maxCardForPlayer = 5

type Game struct {
	board      []Field
	deck       Deck
	usedDeck   Deck
	players    []Player
	round      int
	playerTurn int
	isEnd      bool
	winer      int
}

func (game *Game) GetBoard() []Field {
	return game.board
}
func (game *Game) GetPlayerTurn() int {
	return game.playerTurn + 1
}

func generatePlayers(numberOfPlayers int) []Player {
	players := make([]Player, numberOfPlayers)
	for i := 0; i < numberOfPlayers; i++ {
		players[i] = Player{Color: Colors[i]} ///TODO shuffle the colors
	}
	return players
}
func (game *Game) dealTheCards() {
	for j := 0; j < maxCardForPlayer; j++ {
		for i := range game.players {
			card, _ := game.deck.GetCardFromDeck()
			game.players[i].Cards = append(game.players[i].Cards, card)
		}
	}
}

func (game *Game) getPlayerCards(playerNumber int) ([]Card, error) {

	if playerNumber > len(game.players) {
		return nil, ErrInvalidPlayerNumber
	}
	playerNumber = playerNumber - 1
	if playerNumber < 0 {
		return nil, ErrInvalidPlayerNumber
	}
	return game.players[playerNumber].Cards, nil
}

func (game *Game) playCard(c Card, color Color) (err error, winingPlayer int) {
	player := game.players[game.playerTurn]
	if err := game.checkIfCardAndColorIsValid(c, color); err != nil {
		return err, -1
	}
	if c.typ == LastOne && c.color == Default && color == Default {
		colors := findLastOnePawns(game.board)
		if len(colors) != 1 {
			return ErrPickTheColor, -1
		}
		c.color = Colors[0]
	}
	player.Cards = removeCard(player.Cards, c)
	col := c.color
	if c.color == Default {
		col = color
	}
	b, err := MovePawn(game.board, col, c.move)
	if err != nil {
		return err, -1
	}
	game.board = b

	endGame, color := CheckIfGameOver(game.board)
	if endGame {
		for i, p := range game.players {
			if p.Color == color {
				return nil, i + 1
			}
		}
	}
	newCard, err := game.deck.GetCardFromDeck()
	game.usedDeck = append(game.usedDeck, c)
	if len(game.deck) == 0 {
		game.deck = game.usedDeck
		game.usedDeck = Deck{}
	}
	if err != nil {
		return err, -1
	}
	player.Cards = append(player.Cards, newCard)
	game.players[game.playerTurn] = player
	game.playerTurn = (game.playerTurn) + 1
	if game.playerTurn >= len(game.players) {
		game.playerTurn = 0
	}
	return nil, -1
}
func findLastOnePawns([]Field) []Color {
	for _, f := range []Field{} {
		if len(f.Pawns) > 0 {
			return f.Pawns
		}
	}
	return Colors
}
func (game *Game) checkIfCardAndColorIsValid(card Card, color Color) error {
	player := game.players[game.playerTurn]
	if !checkIfExist(player.Cards, card) {
		return ErrInvalidCard
	}
	if card.typ == Normal && card.color == Default && color == Default {
		return ErrInvalidCard
	}
	colors := findLastOnePawns(game.board)
	if card.typ == LastOne && card.color != Default {
		for _, c := range colors {
			if c == card.color {
				return nil
			}
		}
		return ErrInvalidCard
	}
	return nil
}
func findCard(symbol Symbol) (Card, error) {
	log.Println(symbol)
	for _, card := range DefaultDeck {
		if card.Symbol == symbol {
			return card, nil
		}
	}
	return Card{}, ErrInvalidCard
}

func getColor(text string) Color {
	switch text {
	case "red":
		return Red
	case "blue":
		return Blue
	case "green":
		return Green
	case "yellow":
		return Yellow
	case "purpule":
		return Purpule
	default:
		return Default
	}
}

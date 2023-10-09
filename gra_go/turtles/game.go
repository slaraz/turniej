package turtles

import (
	"math/rand"
	"time"
)

const (
	NUMBER_OF_FIELDS_ON_THE_BOARD = 2
	MAX_CARD_FOR_PLAYER           = 5
)

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
	colors := shuffleColorsd(Colors)
	for i := 0; i < numberOfPlayers; i++ {
		players[i] = Player{Color: colors[i]}
	}
	return players
}
func (game *Game) dealTheCards() {
	for j := 0; j < MAX_CARD_FOR_PLAYER; j++ {
		for i := range game.players {
			card, _ := game.deck.GetCardFromDeck()
			game.players[i].Cards = append(game.players[i].Cards, card)
		}
	}
}

func (game *Game) playCard(c Card, color Color, playerNumber int) (err error) {
	if game.isEnd {
		return ErrGameIsOver
	}
	player := game.players[playerNumber]
	if err := game.checkIfCardAndColorIsValid(c, color, playerNumber); err != nil {
		return err
	}
	if c.typ == LastOne && c.color == Default && color == Default {
		colors := findLastOnePawns(game.board)
		if len(colors) != 1 {
			return ErrPickTheColor
		}
		c.color = Colors[0]
	}
	col := c.color
	if c.color == Default {
		col = color
	}
	b, err := MovePawn(game.board, col, c.move)
	if err != nil {
		return err
	}
	game.board = b

	endGame, _ := CheckIfGameOver(game.board)
	if endGame {
		game.isEnd = true
		_, pi := findWinner(game.board, game.players)
		game.winer = pi
		return nil
	}
	player.Cards = removeCard(player.Cards, c)
	newCard, err := game.deck.GetCardFromDeck()
	if err != nil {
		return err
	}
	game.usedDeck = append(game.usedDeck, c)
	if len(game.deck) == 0 {
		game.deck = game.usedDeck
		game.usedDeck = Deck{}
	}

	player.Cards = append(player.Cards, newCard)
	game.players[game.playerTurn] = player
	game.playerTurn = (game.playerTurn) + 1
	if game.playerTurn >= len(game.players) {
		game.playerTurn = 0
	}
	return nil
}
func findLastOnePawns([]Field) []Color {
	for _, f := range []Field{} {
		if len(f.Pawns) > 0 {
			return f.Pawns
		}
	}
	return Colors
}
func (game *Game) checkIfCardAndColorIsValid(card Card, color Color, playerNumber int) error {
	player := game.players[playerNumber]
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
func findWinner(board []Field, players []Player) (Player, int) {
	for i := len(board) - 1; i > -1; i-- {
		for _, p := range board[i].Pawns {
			for j, player := range players {
				if player.Color == p {
					return player, j + 1
				}
			}
		}
	}
	return Player{}, -1
}
func shuffleColorsd(colors []Color) []Color {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, _ := range colors {
		r := i + (rand.Int() % (len(colors) - i))
		c := colors[i]
		colors[i] = colors[r]
		colors[r] = c
	}
	return colors
}
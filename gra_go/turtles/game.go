package turtles

import (
	"math/rand"
	"time"
)

type UsedCard struct {
	CardSymbol string `json:"cardSymbol"`
	Player     int    `json:"player"`
}

const (
	NUMBER_OF_FIELDS_ON_THE_BOARD = 10
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
	UsedCards  []UsedCard
}

func (game *Game) GetBoard() []Field {
	return game.board
}
func (game *Game) GetPlayerTurn() int {
	return game.playerTurn + 1
}

func generatePlayers(numberOfPlayers int) []Player {
	players := make([]Player, numberOfPlayers)
	colors := shuffleColors(Colors)
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
func (game *Game) removePlayerFromGame(playerNumber int) error {
	if playerNumber > len(game.players) || playerNumber < 0 {
		return ErrInvalidPlayerNumber
	}
	game.players[playerNumber].Color = Default
	return nil
}
func (game *Game) playCard(c Card, color Color, playerNumber int) (err error) {
	if game.isEnd {
		return ErrGameIsOver
	}
	player := game.players[playerNumber]
	if err := game.checkIfCardAndColorIsValid(c, color, playerNumber); err != nil {
		return err
	}
	col := c.color
	if c.typ == LastOne && c.color == Default {
		colors := findLastOnePawns(game.board)
		if len(colors) != 1 && color == Default {
			return ErrPickTheColor
		}
		if !checkIdColorValidDoDL(colors, color) {
			return ErrInvalidColor
		}
		col = Colors[0]
	}

	if col == Default {
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
	game.UsedCards = append(game.UsedCards, UsedCard{CardSymbol: string(c.Symbol), Player: playerNumber})
	if len(game.deck) == 0 {
		game.deck = game.usedDeck
		game.usedDeck = Deck{}
	}

	player.Cards = append(player.Cards, newCard)
	game.players[playerNumber] = player
	game.playerTurn = game.getPlayerTurn(playerNumber)
	if game.playerTurn == -1 {
		return ErrNoPlayers
	}

	return nil
}

func (game *Game) getPlayerTurn(currentPlayer int) int {
	newPlayer := currentPlayer + 1
	a := 0
	for checkIfCorrectPlayerNumber(game.players, currentPlayer) {
		if newPlayer >= len(game.players) {
			newPlayer = 0
		}
		if checkIfCorrectPlayerNumber(game.players, newPlayer) {
			return newPlayer
		}
		newPlayer++
		if a > 10 { //as if there was no player capable of playing
			return -1
		}
		a++
	}
	return -1
}

func checkIfCorrectPlayerNumber(players []Player, playerNumber int) bool {
	if playerNumber > len(players) || playerNumber < 0 {
		return false
	}
	if players[playerNumber].Color == Default {
		return false
	}
	return true
}
func findLastOnePawns(fields []Field) []Color {
	for _, f := range fields {
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
	co := card.color
	if card.color == Default {
		co = color
	}
	if !checkIfTurtleIsOnTheBoard(game.board, co) && card.move < 0 {
		return ErrInvalidCard
	}
	return nil
}
func checkIfTurtleIsOnTheBoard(board []Field, color Color) bool {
	for _, f := range board {
		for _, c := range f.Pawns {
			if c == color {
				return true
			}
		}
	}
	return false
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
	case "purple":
		return Purple
	default:
		return Default
	}
}
func findWinner(board []Field, players []Player) (Player, int) {
	for i := len(board) - 1; i > -1; i-- {
		for p := len(board[i].Pawns) - 1; p > -1; p-- {
			for j, player := range players {
				if player.Color == board[i].Pawns[p] {
					return player, j + 1
				}
			}
		}
	}
	return Player{}, -1
}
func shuffleColors(colors []Color) []Color {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i, _ := range colors {
		r := i + (rand.Int() % (len(colors) - i))
		c := colors[i]
		colors[i] = colors[r]
		colors[r] = c
	}
	return colors
}

func checkIdColorValidDoDL(colors []Color, color Color) bool {
	for _, c := range colors {
		if c == color {
			return true
		}
	}
	return false
}

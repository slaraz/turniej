package turtles

var maxCardForPlayer = 5

type Game struct {
	board      []Field
	deck       Deck
	usedDeck   Deck
	players    []Player
	round      int
	playerTurn int
}

func (game *Game) GetBoard() []Field {
	return game.board
}
func (game *Game) GetPlayerTurn() int {
	return game.playerTurn + 1
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

func (game *Game) GetPlayerCards(playerNumber int) ([]Card, error) {

	if playerNumber > len(game.players) {
		return nil, ErrInvalidPlayerNumber
	}
	playerNumber = playerNumber - 1
	if playerNumber < 0 {
		return nil, ErrInvalidPlayerNumber
	}
	return game.players[playerNumber].Cards, nil
}

func (game *Game) PlayCard(c Card, color Color) (err error, winingPlayer int) {
	player := game.players[game.playerTurn]
	if err := game.CheckIfCardAndColorIsValid(c, color); err != nil {
		return err, -1
	}
	if c.Type == LastOne && c.Color == Default && color == Default {
		colors := findLastOnePawns(game.board)
		if len(colors) != 1 {
			return ErrPickTheColor, -1
		}
		c.Color = Colors[0]
	}
	player.Cards = removeCard(player.Cards, c)
	col := c.Color
	if c.Color == Default {
		col = color
	}
	b, err := MovePawn(game.board, col, c.Move)
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
func (game *Game) CheckIfCardAndColorIsValid(card Card, color Color) error {
	player := game.players[game.playerTurn]
	if !checkIfExist(player.Cards, card) {
		return ErrInvalidCard
	}
	if card.Type == Normal && card.Color == Default && color == Default {
		return ErrInvalidCard
	}
	colors := findLastOnePawns(game.board)
	if card.Type == LastOne && card.Color != Default {
		for _, c := range colors {
			if c == card.Color {
				return nil
			}
		}
		return ErrInvalidCard
	}
	return nil
}

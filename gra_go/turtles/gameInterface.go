package turtles

import (
	"log"
	"strings"

	"github.com/slaraz/turniej/gra_go/proto"
)

type Move struct {
	CardSymbol string `json:"cardSymbol"`
	Color      string `json:"color"`
}

type GameStatus struct {
	Board       []Field    `json:"board"`
	Cards       []Card     `json:"cards"`
	Winer       int        `json:"winer"`
	IsEnd       bool       `json:"isEnd"`
	TurtleColor Color      `json:"turtleColor"`
	UsedCards   []UsedCard `json:"usedCards"`
}

// GetGameStatus - return game status for player
// playerNumber starts from 1
func (game *Game) GetGameStatus(playerNumber int) (*proto.StanGry, error) {
	if playerNumber-1 > len(game.players) || playerNumber-1 < 0 {
		return nil, ErrInvalidPlayerNumber
	}
	status := GameStatus{
		Board:       game.board,
		Cards:       game.players[playerNumber-1].Cards,
		TurtleColor: game.players[playerNumber-1].Color,
		Winer:       game.winer, //IF WINER IS -1 THEN NO WINER
		IsEnd:       game.isEnd,
		UsedCards:   game.UsedCards,
	}
	log.Printf("-----> GetGameStatus: playerNumber: %d, status: %+v", playerNumber, status)
	stat := mapGameStatus(&status)
	log.Printf("-----> StanGry: %+v", stat)
	return stat, nil
}

func mapGameStatus(status *GameStatus) *proto.StanGry {
	return &proto.StanGry{
		TwojKolor:  mapKolor(status.TurtleColor),
		TwojeKarty: mapCards(status.Cards),
		Plansza:    mapBoard(status.Board),
		CzyKoniec:  status.IsEnd,
		KtoWygral:  int32(status.Winer),
	}
}
func mapCards(cards []Card) []proto.Karta {
	karty := []proto.Karta{}
	for _, c := range cards {
		karty = append(karty, mapKarta(c))
	}
	return karty
}
func mapBoard(board []Field) []*proto.Pole {
	pola := []*proto.Pole{}
	for _, b := range board {
		pole := &proto.Pole{
			Zolwie: []proto.KolorZolwia{},
		}
		for _, t := range b.Pawns {
			pole.Zolwie = append(pole.Zolwie, mapKolor(t))
		}
		pola = append(pola, pole)
	}
	return pola
}

func mapKarta(card Card) proto.Karta {
	karta := proto.Karta(proto.Karta_value[string(card.Symbol)])
	return karta
}
func mapKolor(color Color) proto.KolorZolwia {
	kolorZolwia := proto.KolorZolwia(proto.KolorZolwia_value[strings.ToUpper(string(color))])
	return kolorZolwia
}

// Move - player move
func (game *Game) Move(kolor proto.KolorZolwia, cardSymbol proto.Karta, playerNumber int) error {
	playerNumber = playerNumber - 1
	if playerNumber > len(game.players) || playerNumber < 0 {
		return ErrInvalidPlayerNumber
	}

	card, err := findCard(Symbol(proto.Karta_name[int32(cardSymbol)]))
	if err != nil {
		return err
	}
	color := getColor(strings.ToLower(proto.KolorZolwia_name[int32(kolor)]))
	err = game.playCard(card, color, playerNumber)

	return err
}

// CreateNewGame - create new game
func CreateNewGame(numberOfPlayers int) *Game {
	game := &Game{
		board:      CreateGameBoard("a"),
		deck:       CreateGameDeck("a"),
		players:    generatePlayers(numberOfPlayers),
		round:      0,
		playerTurn: 0,
	}
	game.dealTheCards()
	return game
}

func (game *Game) RemovePlayer(playerNumber int) error {
	playerNumber = playerNumber - 1
	return game.removePlayerFromGame(playerNumber)
}

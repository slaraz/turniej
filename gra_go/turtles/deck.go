package turtles

import (
	"errors"
	"time"

	"math/rand"
)

type Symbol string
type Deck []Card

var gameDecks = map[string][]Card{}

type Type string

type Card struct {
	id     int
	color  Color
	move   int
	typ    Type
	Symbol Symbol `json:"symbol"`
}

func shuffleDefaultDeck() Deck {
	cards := make([]Card, len(DefaultDeck))
	copy(cards, DefaultDeck)
	return shuffle(cards)
}

func shuffle(deck Deck) Deck {
	for i, _ := range deck {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))

		r := i + (rand.Int() % (len(deck) - i))
		c := deck[i]
		deck[i] = deck[r]
		deck[r] = c

	}
	return deck
}

func CreateGameDeck(uuid string) Deck {
	gameDecks[uuid] = shuffleDefaultDeck()
	return gameDecks[uuid]
}

func (deck *Deck) GetCardFromDeck() (Card, error) {
	if len(*deck) == 0 {
		return Card{}, errors.New("no cards in deck")
	}
	card := (*deck)[0]
	*deck = (*deck)[1:]
	return card, nil
}
func DeleteDeck(uuid string) {
	delete(gameDecks, uuid)
}

var Normal = Type("normal")
var LastOne = Type("lastone")
var R1 = Symbol("R1")
var R2 = Symbol("R2")
var R3 = Symbol("R1B")
var G1 = Symbol("G1")
var G2 = Symbol("G2")
var G3 = Symbol("G1B")
var B1 = Symbol("B1")
var B2 = Symbol("B2")
var B3 = Symbol("B1B")
var Y1 = Symbol("Y1")
var Y2 = Symbol("Y2")
var Y3 = Symbol("Y1B")
var P1 = Symbol("P1")
var P2 = Symbol("P2")
var P3 = Symbol("P1B")
var D1 = Symbol("A1")
var D3 = Symbol("A1B")
var DL1 = Symbol("L1")
var DL2 = Symbol("L2")

var DefaultDeck = []Card{
	{1, Red, 1, Normal, R1},
	{2, Red, 1, Normal, R1},
	{3, Red, 1, Normal, R1},
	{4, Red, 1, Normal, R1},
	{5, Red, 1, Normal, R1},
	{6, Blue, 1, Normal, B1},
	{7, Blue, 1, Normal, B1},
	{8, Blue, 1, Normal, B1},
	{9, Blue, 1, Normal, B1},
	{10, Blue, 1, Normal, B1},
	{11, Green, 1, Normal, G1},
	{12, Green, 1, Normal, G1},
	{13, Green, 1, Normal, G1},
	{14, Green, 1, Normal, G1},
	{15, Green, 1, Normal, G1},
	{16, Yellow, 1, Normal, Y1},
	{17, Yellow, 1, Normal, Y1},
	{18, Yellow, 1, Normal, Y1},
	{19, Yellow, 1, Normal, Y1},
	{20, Yellow, 1, Normal, Y1},
	{21, Purple, 1, Normal, P1},
	{22, Purple, 1, Normal, P1},
	{23, Purple, 1, Normal, P1},
	{24, Purple, 1, Normal, P1},
	{25, Purple, 1, Normal, P1},
	{26, Red, 2, Normal, R2},
	{27, Yellow, 2, Normal, Y2},
	{28, Green, 2, Normal, G2},
	{29, Blue, 2, Normal, B2},
	{30, Purple, 2, Normal, P2},
	{31, Red, -1, Normal, R3},
	{32, Red, -1, Normal, R3},
	{33, Yellow, -1, Normal, Y3},
	{34, Yellow, -1, Normal, Y3},
	{35, Green, -1, Normal, G3},
	{36, Green, -1, Normal, G3},
	{37, Blue, -1, Normal, B3},
	{38, Blue, -1, Normal, B3},
	{39, Purple, -1, Normal, P3},
	{40, Purple, -1, Normal, P3},
	{41, Default, 1, Normal, D1},
	{42, Default, 1, Normal, D1},
	{43, Default, 1, Normal, D1},
	{44, Default, 1, Normal, D1},
	{45, Default, 1, Normal, D1},
	{46, Default, -1, Normal, D3},
	{47, Default, -1, Normal, D3},
	{48, Default, 1, LastOne, DL1},
	{49, Default, 1, LastOne, DL1},
	{50, Default, 1, LastOne, DL1},
	{51, Default, 2, LastOne, DL2},
	{52, Default, 2, LastOne, DL2},
}

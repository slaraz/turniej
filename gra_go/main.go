package main

import (
	"bufio"
	game "hackaton/turtles"
	"strconv"
	"strings"

	"log"
	"os"
)

// TODO
// 2. Fix moving backwords to the start
// 3. Add wining condition when pawn is difrent then players color
// 4. Decide what to do if the pawn jump over other the last field
// 5. Create card shortcuts
func main() {
	log.Println("Hello World")

	reader := bufio.NewReader(os.Stdin)
	//log.Println("Enter text: ")
	game := game.CreateNewGame(2)
	log.Println(game.GetBoard())
	log.Println(game.GetPlayerTurn())
	log.Println(game.GetPlayerCards(1))
	log.Println(game.GetPlayerCards(2))
	log.Println(game.GetPlayerCards(3))
	log.Println(game.GetPlayerCards(4))

	winer := 0

	for winer < 1 {
		log.Println("Enter text: color,move")
		text, _ := reader.ReadString('\n')
		card, color, err := getCardFromText(text, game)
		if err != nil {
			log.Println(err)
			continue
		}
		err, winer = game.PlayCard(card, color)
		if winer > 0 {
			break
		}
		if err != nil {
			log.Println(err)
			continue

		}
		log.Println(game.GetBoard())
		log.Printf("Player: %d", game.GetPlayerTurn())
		log.Println(game.GetPlayerCards(game.GetPlayerTurn()))
	}
	log.Println("Winer is player: ", winer)

}
func getCardFromText(text string, g game.Game) (game.Card, game.Color, error) {

	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "", -1)
	c := strings.Split(text, ",")
	if len(c) < 2 {
		c = strings.Split(text, " ")
	}
	if len(c) < 2 {
		return game.Card{}, game.Default, game.ErrInvalidCard
	}

	a, err := strconv.Atoi(c[0])
	color := getColor(c[1])
	if err != nil {
		return game.Card{}, game.Default, err
	}

	playerCards, err := g.GetPlayerCards(g.GetPlayerTurn())
	if err != nil {
		return game.Card{}, game.Default, err
	}
	return playerCards[a], color, nil
}
func getColor(text string) game.Color {
	switch text {
	case "red":
		return game.Red
	case "blue":
		return game.Blue
	case "green":
		return game.Green
	case "yellow":
		return game.Yellow
	case "purpule":
		return game.Purpule
	default:
		return game.Default
	}
}

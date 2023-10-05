package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"

	"log"
	"os"

	"github.com/slaraz/turniej/gra_go/turtles"
)

// TODO

// 9. change number of fields
// 10. server picks the player
func main() {
	log.Println("Hello World")

	reader := bufio.NewReader(os.Stdin)
	//log.Println("Enter text: ")
	game := turtles.CreateNewGame(2)

	winer := 0
	isEnd := false
	for isEnd == false {
		log.Println("Enter text: symbol,color")
		log.Printf("Player: %d", game.GetPlayerTurn())
		res1, _ := game.GetGameStatus(game.GetPlayerTurn())
		log.Println(res1)
		text, _ := reader.ReadString('\n')
		move, err := getCardFromText(text)
		if err != nil {
			log.Println(err)
			continue
		}
		str, _ := json.Marshal(move)
		err = game.Move(string(str))
		if err != nil {
			log.Println(err)
			continue
		}

		res, _ := game.GetGameStatus(game.GetPlayerTurn())
		st := turtles.GameStatus{}
		json.Unmarshal([]byte(res), &st)
		isEnd = st.IsEnd
		winer = st.Winer

	}
	log.Println("Winer is player: ", winer)

}
func getCardFromText(text string) (turtles.Move, error) {

	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "", -1)
	c := strings.Split(text, ",")
	if len(c) < 2 {
		c = strings.Split(text, " ")
	}
	if len(c) < 2 {
		return turtles.Move{}, fmt.Errorf("invalid input")
	}

	a := c[0]
	color := getColor(strings.ToLower(c[1]))
	m := turtles.Move{CardSymbol: a, Color: string(color)}

	return m, nil

}
func getColor(text string) turtles.Color {
	switch text {
	case "red":
		return turtles.Red
	case "blue":
		return turtles.Blue
	case "green":
		return turtles.Green
	case "yellow":
		return turtles.Yellow
	case "purpule":
		return turtles.Purpule
	default:
		return turtles.Default
	}
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"log"
	"os"

	"github.com/slaraz/turniej/gra_go/turtles"
)

// TODO
// 9. change number of fieldsD

func main() {
	log.Println("Hello World")

	reader := bufio.NewReader(os.Stdin)
	//log.Println("Enter text: ")
	game := turtles.CreateNewGame(3)

	winer := 0
	isEnd := false
	for !isEnd {
		log.Println("Enter text: symbol,playerNumber,color")
		log.Printf("Player: %d", game.GetPlayerTurn())
		res1, _ := game.GetGameStatus(game.GetPlayerTurn())
		log.Println(res1)
		text, _ := reader.ReadString('\n')
		move, playerNumber, err := getCardFromText(text)
		if err != nil {
			log.Println(err)
			continue
		}
		str, _ := json.Marshal(move)
		err = game.Move(string(str), playerNumber)
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
func getCardFromText(text string) (turtles.Move, int, error) {

	text = strings.Replace(text, "\r", "", -1)
	text = strings.Replace(text, "\n", "", -1)
	c := strings.Split(text, ",")
	if len(c) < 3 {
		c = strings.Split(text, " ")
	}
	if len(c) < 3 {
		return turtles.Move{}, 0, fmt.Errorf("invalid input")
	}
	a := c[0]
	color := getColor(strings.ToLower(c[2]))
	m := turtles.Move{CardSymbol: a, Color: string(color)}
	playerNumber, err := strconv.Atoi(c[1])
	return m, playerNumber, err
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

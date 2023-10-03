package turtles

type Player struct {
	Color Color
	Cards []Card
}

func checkIfExist(cards []Card, card Card) bool {
	for _, c := range cards {
		if c.Color == card.Color && c.Move == card.Move {
			return true
		}
	}
	return false
}
func removeCard(cards []Card, c Card) []Card {
	for i, card := range cards {
		if card.Color == c.Color && card.Move == c.Move {
			return append(cards[:i], cards[i+1:]...)
		}
	}
	return cards
}

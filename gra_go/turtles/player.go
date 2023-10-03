package turtles

type Player struct {
	Color Color
	Cards []Card
}

func checkIfExist(cards []Card, card Card) bool {
	for _, c := range cards {
		if c.color == card.color && c.move == card.move {
			return true
		}
	}
	return false
}
func removeCard(cards []Card, c Card) []Card {
	for i, card := range cards {
		if card.color == c.color && card.move == c.move {
			return append(cards[:i], cards[i+1:]...)
		}
	}
	return cards
}

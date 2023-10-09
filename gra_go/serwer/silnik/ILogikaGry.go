package silnik

import (
	"github.com/slaraz/turniej/gra_go/proto"
	"github.com/slaraz/turniej/gra_go/turtles"
)

type ILogikaGry interface {
	GetGameStatus(playerNumber int) (*proto.StanGry, error)
	Move(playerColor proto.KolorZolwia, cardSymbol proto.Karta, playerNumber int) error
}

func getLogikaGry(liczbaGraczy int) ILogikaGry {
	return turtles.CreateNewGame(liczbaGraczy)
}

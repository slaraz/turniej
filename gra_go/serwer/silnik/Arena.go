package silnik

import (
	"fmt"

	"github.com/slaraz/turniej/gra_go/proto"
)

const (
	DLUGOSC_ID       = 3
	DLUGOSC_GRACZ_ID = 10
)

type ArenaGry struct {
	aktywneGry map[string]*gra
}

func NowaArena() *ArenaGry {
	arena := &ArenaGry{
		aktywneGry: map[string]*gra{},
	}
	return arena
}

func (sil *ArenaGry) NowaGra(iluGraczy int) (string, error) {
	graId := ""
	for {
		graId = generujLosoweId(DLUGOSC_ID)
		// czy jest takie id?
		if _, ok := sil.aktywneGry[graId]; !ok {
			// nie ma, bierzemy
			break
		}
	}
	taGra := nowaGra(graId, iluGraczy)
	sil.aktywneGry[graId] = taGra

	return graId, nil
}

func (arena *ArenaGry) DodajGraczaDoGry(graId string, wizytowka *proto.WizytowkaGracza) (string, error) {
	gra, ok := arena.aktywneGry[graId]
	if !ok {
		return "", fmt.Errorf("brak aktywnej gry %q", graId)
	}
	graczId, err := gra.dodajGracza(wizytowka)
	if err != nil {
		return "", err
	}
	return graczId, nil
}

func (arena *ArenaGry) StanGry(graId, graczId string) (*proto.StanGry, error) {
	gra, ok := arena.aktywneGry[graId]
	if !ok {
		return nil, fmt.Errorf("brak aktywnej gry %q", graId)
	}
	return gra.stanGry(graczId)
}

func (arena *ArenaGry) RuchGracza(ruch *proto.RuchGracza) error {
	gra, ok := arena.aktywneGry[ruch.GraId]
	if !ok {
		return fmt.Errorf("brak aktywnej gry %q", ruch.GraId)
	}

	return gra.ruchGracza(ruch)
}

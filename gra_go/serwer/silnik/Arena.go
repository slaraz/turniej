package silnik

import (
	"fmt"
	"log"

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
	go arena.arenaFlow()

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

	go func() {
		err := taGra.przebiegRozgrywki()
		if err != nil {
			log.Println("błąd przebiegRozgrywki:", err)
		}
		// TODO: współbierznie zrobić
		delete(sil.aktywneGry, graId)
	}()

	return graId, nil
}

func (arena *ArenaGry) DodajGraczaDoGry(graId string, wizytowka *proto.WizytowkaGracza) (string, error) {
	gra, ok := arena.aktywneGry[graId]
	if !ok {
		return "", fmt.Errorf("brak aktywnej gry %q", graId)
	}
	graczId, err := gra.DodajGracza(wizytowka)
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
	arena.kanGetGra <- struct {ruch.GraId; kanGra}
	gra, err := <-kanGra
	gra := <-arena.kanGry
//	gra, ok := arena.aktywneGry[ruch.GraId]
	if err!= nil {
		return fmt.Errorf("brak gry %q", ruch.GraId)
	}

	gra.kanRuchGracza <- struct {ruch; kanOdp}
	odp := <- kanOdp
	if odp == koniecGry {
		arena.kanKoniecGry <-ruch.GraId
	}
	return gra.ruchGracza(ruch)
}

func (sil *ArenaGry) arenaFlow() {

	nowaGra := nowaGra()
	sil.aktywneGry[nowaGra.graId] = nowaGra
	go nowaGra.Rozgrywka()

	for {
	select {
	case graId, kanGra := <- arena.kanGetGra:
		gra, ok := arena.aktywneGry[graId] 
		if !ok {
			kanGra <- struct{ nil; fmt.Errorf("arena.aktywneGry: %q", graId)
			continue
		}
		kanGra <- struct{gra; nil}
	case graId <- kanKoniecGry:
		delete(sil.aktywneGry, graId)
	}
	}

}
}



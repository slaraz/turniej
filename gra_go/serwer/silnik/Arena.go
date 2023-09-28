package silnik

import (
	"fmt"
)

type ArenaGry struct {
	aktywneGry map[string]*gra
	kanNewGra  chan reqNewGra
	kanGetGra  chan reqGetGra
	kanEndGra  chan reqEndGra
}

func NowaArena() *ArenaGry {
	arena := &ArenaGry{
		aktywneGry: map[string]*gra{},
		kanNewGra:  make(chan reqNewGra),
		kanGetGra:  make(chan reqGetGra),
		kanEndGra:  make(chan reqEndGra),
	}
	go arena.arenaFlow()

	return arena
}

type reqNewGra struct {
	iluGraczy int
	kanOdp    chan odpNowaGra
}

type odpNowaGra struct {
	graId string
	err   error
}

func (arena *ArenaGry) NowaGra(iluGraczy int) (string, error) {
	kanOdp := make(chan odpNowaGra)
	arena.kanNewGra <- reqNewGra{
		iluGraczy: iluGraczy,
		kanOdp:    kanOdp,
	}
	odp := <-kanOdp
	return odp.graId, odp.err
}

type reqGetGra struct {
	graID  string
	kanOdp chan odpGetGra
}

type odpGetGra struct {
	gra *gra
	err error
}

func (arena *ArenaGry) GetGra(graID string) (*gra, error) {
	kanOdp := make(chan odpGetGra)
	arena.kanGetGra <- reqGetGra{
		graID:  graID,
		kanOdp: kanOdp,
	}
	odp := <-kanOdp
	return odp.gra, odp.err
}

type reqEndGra struct {
	graID  string
	kanOdp chan odpEndGra
}

type odpEndGra struct {
	err error
}

func (arena *ArenaGry) KoniecGry(graID string) error {
	kanOdp := make(chan odpEndGra)
	arena.kanEndGra <- reqEndGra{
		graID:  graID,
		kanOdp: kanOdp,
	}
	odp := <-kanOdp
	return odp.err
}

func (arena *ArenaGry) arenaFlow() {
	for {
		select {

		case req := <-arena.kanNewGra:
			//TODO: zrobić ograniczenie liczby gier per serwer
			graId := arena.getNowaGraID()
			arena.aktywneGry[graId] = nowaGra(graId, req.iluGraczy)
			req.kanOdp <- odpNowaGra{
				graId: graId,
			}

		case req := <-arena.kanGetGra:
			odp := odpGetGra{}
			gra, ok := arena.aktywneGry[req.graID]
			if !ok {
				odp.err = fmt.Errorf("brak gry arena.aktywneGry[%q]", req.graID)
			} else {
				odp.gra = gra
			}
			req.kanOdp <- odp

		case req := <-arena.kanEndGra:
			delete(arena.aktywneGry, req.graID)
			req.kanOdp <- odpEndGra{}
		}
	}
}

// func (arena *ArenaGry) DodajGraczaDoGry(graId string, wizytowka *proto.WizytowkaGracza) (string, error) {

// 	gra, ok := arena.aktywneGry[graId]
// 	if !ok {
// 		return "", fmt.Errorf("brak aktywnej gry %q", graId)
// 	}
// 	graczId, err := gra.DodajGracza(wizytowka)
// 	if err != nil {
// 		return "", err
// 	}
// 	return graczId, nil
// }

// func (arena *ArenaGry) RuchGracza(ruch *proto.RuchGracza) error {
// 	kanGra := make(chan struct {
// 		*gra
// 		error
// 	})
// 	arena.kanGetGra <- struct {
// 		ruch.GraId
// 		kanGra
// 	}
// 	gra, err := <-kanGra
// 	if err != nil {
// 		return fmt.Errorf("RuchGracza gra[%q]: %v", ruch.GraId, err)
// 	}

// 	gra.kanRuchGracza <- struct {
// 		ruch
// 		kanOdp
// 	}
// 	odp := <-kanOdp
// 	if odp == koniecGry {
// 		arena.kanKoniecGry <- ruch.GraId
// 	}
// 	return gra.ruchGracza(ruch)
// }

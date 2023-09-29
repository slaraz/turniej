package silnik

import (
	"fmt"
)

type ArenaGry struct {
	aktywneGry   map[string]*gra
	kanNewGra    chan reqNewGra
	kanGetGra    chan reqGetGra
	kanKoniecGry chan reqKoniecGry
}

func NowaArena() *ArenaGry {
	arena := &ArenaGry{
		aktywneGry:   map[string]*gra{},
		kanNewGra:    make(chan reqNewGra),
		kanGetGra:    make(chan reqGetGra),
		kanKoniecGry: make(chan reqKoniecGry),
	}
	go arena.arenaFlow()

	return arena
}

type reqNewGra struct {
	iluGraczy int
	kanOdp    chan odpNowaGra
}

type odpNowaGra struct {
	graID string
	err   error
}

func (arena *ArenaGry) NowaGra(iluGraczy int) (string, error) {
	kanOdp := make(chan odpNowaGra)
	arena.kanNewGra <- reqNewGra{
		iluGraczy: iluGraczy,
		kanOdp:    kanOdp,
	}
	odp := <-kanOdp
	return odp.graID, odp.err
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

type reqKoniecGry struct {
	graID string
	err   error
}

type odpEndGra struct {
	err error
}

// func (arena *ArenaGry) KoniecGry(graID string) error {
// 	kanOdp := make(chan odpEndGra)
// 	arena.kanEndGra <- reqKoniecGry{
// 		graID:  graID,
// 		kanOdp: kanOdp,
// 	}
// 	odp := <-kanOdp
// 	return odp.err
// }

func (arena *ArenaGry) arenaFlow() {
	for {
		select {

		case req := <-arena.kanNewGra:
			odp := odpNowaGra{}
			//TODO: zrobić ograniczenie liczby gier per serwer
			// generuję unikalny ID gry
			graID := arena.getNowaGraID()
			// uruchamiam nową grę
			gra, err := nowaGra(graID, req.iluGraczy, arena.kanKoniecGry)
			if err != nil {
				odp.err = err
			} else {
				// dokładam nową grę do aktywnych na arenie
				arena.aktywneGry[graID] = gra
				odp.graID = graID
			}
			req.kanOdp <- odp

		case req := <-arena.kanGetGra:
			odp := odpGetGra{}
			gra, ok := arena.aktywneGry[req.graID]
			if !ok {
				odp.err = fmt.Errorf("brak gry arena.aktywneGry[%q]", req.graID)
			} else {
				odp.gra = gra
			}
			req.kanOdp <- odp

		case req := <-arena.kanKoniecGry:
			fmt.Printf("koniec gry %q: %v\n", req.graID, req.err)
			delete(arena.aktywneGry, req.graID)
		}
	}
}

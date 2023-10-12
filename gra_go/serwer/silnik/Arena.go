package silnik

import (
	"fmt"
	"log"
)

type Arena struct {
	aktywneGry   map[string]*gra
	kanNewGra    chan reqNewGra
	kanGetGra    chan reqGetGra
	kanKoniecGry chan reqKoniecGry
}

func NowaArena() *Arena {
	arena := &Arena{
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
	graID     string
	kanOdp    chan odpNowaGra
}

type odpNowaGra struct {
	graID string
	err   error
}

func (arena *Arena) NowaGra(iluGraczy int, graID string) (string, error) {
	kanOdp := make(chan odpNowaGra)
	arena.kanNewGra <- reqNewGra{
		iluGraczy: iluGraczy,
		graID:     graID,
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

func (arena *Arena) GetGra(graID string) (*gra, error) {
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

func (arena *Arena) arenaFlow() {
	for {
		select {

		case req := <-arena.kanNewGra:
			odp := odpNowaGra{}
			// generuję unikalny ID gry
			var graID string
			if req.graID == "" {
				// nie podano graID, losujemy nowy
				graID = arena.getNowaGraID()
			} else {
				// podano graID
				if _, ok := arena.aktywneGry[req.graID]; ok {
					// ale takie już jest
					odp.err = fmt.Errorf("gra %q już istnieje", req.graID)
				}
				graID = req.graID
			}
			// jeśli z graID ok
			if odp.err == nil {

				// uruchamiam nową grę
				gra, err := nowaGra(graID, req.iluGraczy, arena.kanKoniecGry)
				if err != nil {
					odp.err = err
				} else {
					// dokładam nową grę do aktywnych na arenie
					arena.aktywneGry[graID] = gra
					odp.graID = graID
				}
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
			log.Printf("%s arenaFlow(): usuwam grę: %v", req.graID, req.err)
			delete(arena.aktywneGry, req.graID)
		}
	}
}

package silnik

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/slaraz/turniej/gra_go/proto"
)

const (
	MIN_LICZBA_GRACZY         = 1
	MAX_LICZBA_GRACZY         = 5
	DOLACZANIE_GRACZY_TIMEOUT = time.Second * 30
	RUCH_GRACZA_TIMEOUT       = time.Second * 1000
	WYSLIJ_STATUS_TIMEOUT     = time.Second * 10
)

type gra struct {
	graID        string
	liczbaGraczy int
	graczeByID   map[string]*gracz
	logika       ILogikaGry

	kanDolaczGracza   chan reqDolaczGracza
	kanArenaKoniecGry chan reqKoniecGry
}

type gracz struct {
	graczID     string
	nazwaGracza string
	kanRuch     chan reqRuchGracza
	kanStan     chan *proto.StanGry
}

func nowaGra(graID string, liczbaGraczy int, kanKoniecGry chan reqKoniecGry) (*gra, error) {
	// ograniczamy liczbę graczy do limitu ustalonego na serwerze
	if liczbaGraczy < MIN_LICZBA_GRACZY {
		return nil, fmt.Errorf("zbyt duża liczba graczy: %d, maksymalna dozwolona liczba graczy to: %d", liczbaGraczy, MAX_LICZBA_GRACZY)
	}
	if liczbaGraczy > MAX_LICZBA_GRACZY {
		return nil, fmt.Errorf("za mała liczba graczy: %d, minimalna dozwolona liczba graczy to: %d", liczbaGraczy, MAX_LICZBA_GRACZY)
	}
	// obiekt gry,
	// do którego odwołują się metody gry
	g := &gra{
		graID:             graID,
		liczbaGraczy:      liczbaGraczy,
		graczeByID:        map[string]*gracz{},
		logika:            getLogikaGry(liczbaGraczy),
		kanDolaczGracza:   make(chan reqDolaczGracza),
		kanArenaKoniecGry: kanKoniecGry,
	}
	log.Printf("%s nowaGra(): utworzono grę dla %d graczy\n", g.graID, g.liczbaGraczy)
	// uruchomienie wątku gry
	go g.przebiegRozgrywki()
	return g, nil
}

type reqDolaczGracza struct {
	nazwaGracza string
	kanOdp      chan odpDolaczGracza
}

type odpDolaczGracza struct {
	gracz *gracz
	err   error
}

func (g *gra) DolaczGracza(nazwaGracza string) (string, error) {
	kanOdp := make(chan odpDolaczGracza)
	g.kanDolaczGracza <- reqDolaczGracza{
		nazwaGracza: nazwaGracza,
		kanOdp:      kanOdp,
	}
	odp := <-kanOdp
	return odp.gracz.graczID, odp.err
}

type reqRuchGracza struct {
	ruch   *proto.RuchGracza
	kanOdp chan odpRuchGracza
}

type odpRuchGracza struct {
	err error
}

func (g *gra) WykonajRuch(graczID string, ruch *proto.RuchGracza) (string, error) {
	gracz, ok := g.graczeByID[graczID]
	if !ok {
		return "", fmt.Errorf("%s WykonajRuch((): nie ma gracza: %q", g.graID, graczID)
	}
	kanOdp := make(chan odpRuchGracza)
	log.Printf("%s WykonajRuch(): %s rząda wykonania ruchu %q %q\n", g.graID, gracz.nazwaGracza, ruch.ZagranaKarta, ruch.KolorWybrany)
	gracz.kanRuch <- reqRuchGracza{
		ruch:   ruch,
		kanOdp: kanOdp,
	}
	odp := <-kanOdp
	if odp.err != nil {
		log.Printf("%s WykonajRuch(): %s wykonał błędny ruch %q: %v\n", g.graID, gracz.nazwaGracza, ruch.ZagranaKarta, odp.err)
	} else {
		log.Printf("%s WykonajRuch(): %s wykonanał ruch %q\n", g.graID, gracz.nazwaGracza, ruch.ZagranaKarta)
	}
	return gracz.graczID, odp.err
}

func (g *gra) StanGry(graczID string) (*proto.StanGry, error) {
	gracz, ok := g.graczeByID[graczID]
	if !ok {
		return nil, fmt.Errorf("%s StanGry(): nie ma gracza: %q", g.graID, graczID)
	}
	log.Printf("%s StanGry(): %s rząda status\n", g.graID, gracz.nazwaGracza)
	stan, ok := <-gracz.kanStan
	// TODO: albo kanał koniec gry
	if !ok {
		return nil, fmt.Errorf("gracz status !ok")
	}
	stan.GraID = g.graID
	stan.GraczID = graczID

	log.Printf("%s StanGry(): %s dostaje: plansza: %v, karty: %v", g.graID, gracz.nazwaGracza, stan.Plansza, stan.TwojeKarty)
	drukujStan(stan)
	return stan, nil
}

func drukujStan(stan *proto.StanGry) {
	stanJSON, _ := json.Marshal(stan)
	log.Printf("stan: %s", stanJSON)
}
func (g *gra) przebiegRozgrywki() {
	log.Printf("%s Rozgrywka0: rozpoczęcie rozgrywki\n", g.graID)

	// dołączanie graczy
	gracze := []*gracz{}
	timeout := time.After(DOLACZANIE_GRACZY_TIMEOUT)
	for i := 1; i <= g.liczbaGraczy; i++ {
		select {
		case req := <-g.kanDolaczGracza:
			log.Printf("%s Rozgrywka1: dolaczanie gracza: %s\n", g.graID, req.nazwaGracza)
			graczID := g.getNowyGraczID()
			nowyGracz := &gracz{
				graczID:     graczID,
				nazwaGracza: req.nazwaGracza,
				kanRuch:     make(chan reqRuchGracza),
				kanStan:     make(chan *proto.StanGry),
			}
			gracze = append(gracze, nowyGracz)
			g.graczeByID[graczID] = nowyGracz
			odp := odpDolaczGracza{
				gracz: nowyGracz,
				err:   nil,
			}
			req.kanOdp <- odp

		case <-timeout:
			// każdy utworzony gracz musi dostać status koniec gry
			for _, gracz := range gracze {
				gracz.kanStan <- &proto.StanGry{
					CzyKoniec: true,
					KtoWygral: -1,
				}
			}
			g.koniec(fmt.Errorf("nie zebrał się komplet graczy w odpowiednim czasie: %v", DOLACZANIE_GRACZY_TIMEOUT))
			return
		}
	}

	// losujemy kolejność graczy
	rand.Shuffle(len(gracze), func(i, j int) {
		gracze[i], gracze[j] = gracze[j], gracze[i]
	})
	kolejnosc := ""
	for i, gracz := range gracze {
		kolejnosc += fmt.Sprintf(" %d. %q", i+1, gracz.nazwaGracza)
	}
	log.Printf("%s Rozgrywka2: kolejność graczy: %s", g.graID, kolejnosc)

	// wykonywanie ruchów
	i := 0
	for {
		ruszajacyGracz := gracze[i]

		// status ---->

		// LOGIKA GRY
		stan, _ := g.logika.GetGameStatus(i + 1)
		// wysyłamy status do gracza
		timeout2 := time.After(WYSLIJ_STATUS_TIMEOUT)
		select {
		case ruszajacyGracz.kanStan <- stan:
		case <-timeout2:
			g.koniec(fmt.Errorf("%s upłynął czas dla gracza: %s", g.graID, ruszajacyGracz.nazwaGracza))
			return
		}
		log.Printf("%s Rozgrywka3: wysłano status dla gracza %q\n", g.graID, ruszajacyGracz.nazwaGracza)

		// koniec gry - rozsyłam statusy do pozostałych graczy
		if stan.CzyKoniec {
			log.Printf("%s Rozgrywka7: rozsyłam statusy", g.graID)
			// jeśli koniec gry wysyłamy status do reszty graczy
			wg := sync.WaitGroup{}
			for i, gr := range gracze {
				if gr.graczID == ruszajacyGracz.graczID {
					continue
				}
				wg.Add(1)
				go func(ii int, ggr *gracz) {
					defer wg.Done()
					// LOGIKA GRY
					stan, _ := g.logika.GetGameStatus(ii + 1)
					// wysyłamy status do gracza
					timeout2 := time.After(WYSLIJ_STATUS_TIMEOUT)
					select {
					case ggr.kanStan <- stan:
					case <-timeout2:

						// TODO: tu dalej pisać jak sobie pójdziesz

						g.koniec(fmt.Errorf("%s upłynął czas dla gracza: %s", g.graID, ggr.nazwaGracza))
						return
					}
					log.Printf("%s Rozgrywka8: wysłano status dla gracza %q\n", g.graID, ggr.nazwaGracza)
				}(i, gr)
			}
			wg.Wait()
			g.koniec(nil)
			return
		}
		log.Printf("%s Rozgrywka4: ruszający gracz %d %q\n", g.graID, i, ruszajacyGracz.nazwaGracza)

		// ruch ---->

		timeout1 := time.After(RUCH_GRACZA_TIMEOUT)
		for {
			select {
			case req := <-ruszajacyGracz.kanRuch:

				// LOGIKA GRY
				err := g.logika.Move(req.ruch.KolorWybrany, req.ruch.ZagranaKarta, i+1)
				if err != nil {
					log.Printf("%s Rozgrywka9: błąd logika.Move(): %v", g.graID, err)
					req.kanOdp <- odpRuchGracza{err: err}
					continue
				}
				req.kanOdp <- odpRuchGracza{}
			case <-timeout1:
				//TODO: rozesłać status kto wygrał, kto przegrał
				g.koniec(fmt.Errorf("%s upłynął czas ruchu dla gracza: %s", g.graID, ruszajacyGracz.nazwaGracza))
				return
			}
			break
		}
		log.Printf("%s Rozgrywka5: wykonano ruch gracza %q\n", g.graID, ruszajacyGracz.nazwaGracza)

		i = g.nastepny(i)
	} //for
}

func (g *gra) koniec(err error) {
	log.Printf("%s KONIEC rozgrywki: %v\n", g.graID, err)
	g.kanArenaKoniecGry <- reqKoniecGry{
		graID: g.graID,
		err:   err,
	}
}

func (g *gra) nastepny(i int) int {
	nast := i + 1
	if nast == g.liczbaGraczy {
		return 0
	} else {
		return nast
	}
}

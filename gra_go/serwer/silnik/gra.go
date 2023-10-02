package silnik

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/slaraz/turniej/logika"
)

const (
	MIN_LICZBA_GRACZY         = 1
	MAX_LICZBA_GRACZY         = 5
	DOLACZANIE_GRACZY_TIMEOUT = time.Second * 300
	RUCH_GRACZA_TIMEOUT       = time.Second * 10
	WYSLIJ_STATUS_TIMEOUT     = time.Second * 10
)

type gra struct {
	graID        string
	liczbaGraczy int
	stolik       map[string]*gracz
	logika       ILogikaGry

	kanDolaczGracza chan reqDolaczGracza
	kanKoniecGry    chan error
}

type ILogikaGry interface {
	Ruch(nrGracza int, ruch string) error
	StanGry(nrGracza int) string
}

type gracz struct {
	graczID     string
	nazwaGracza string
	kanRuch     chan reqRuchGracza
	kanStatus   chan string
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
		graID:           graID,
		liczbaGraczy:    liczbaGraczy,
		stolik:          map[string]*gracz{},
		logika:          &logika.LogikaGry{},
		kanDolaczGracza: make(chan reqDolaczGracza),
		kanKoniecGry:    make(chan error),
	}
	log.Printf("nowaGra: gra %s: utworzono grę dla %d graczy\n", g.graID, g.liczbaGraczy)
	// uruchomienie wątku gry
	go g.przebiegRozgrywki()
	return g, nil
}

type reqDolaczGracza struct {
	nazwaGracza string
	kanOdp      chan odpDolaczGracza
}

type odpDolaczGracza struct {
	graczID string
	err     error
}

func (g *gra) DolaczGracza(nazwaGracza string) (string, error) {
	kanOdp := make(chan odpDolaczGracza)
	g.kanDolaczGracza <- reqDolaczGracza{
		nazwaGracza: nazwaGracza,
		kanOdp:      kanOdp,
	}
	odp := <-kanOdp
	return odp.graczID, odp.err
}

type reqRuchGracza struct {
	ruch   string
	kanOdp chan odpRuchGracza
}

type odpRuchGracza struct {
	err error
}

func (g *gra) WykonajRuch(graczID string, ruch string) error {
	gracz, ok := g.stolik[graczID]
	if !ok {
		return fmt.Errorf("WykonajRuch: nie ma gracza: %q", graczID)
	}
	kanOdp := make(chan odpRuchGracza)
	log.Printf("WykonajRuch: gra %s, gracz %q: rząda wykonania ruchu %q\n", g.graID, gracz.nazwaGracza, ruch)
	gracz.kanRuch <- reqRuchGracza{
		ruch:   ruch,
		kanOdp: kanOdp,
	}
	odp := <-kanOdp
	log.Printf("WykonajRuch: gra %s, gracz %q: wykonano ruchu %q\n", g.graID, gracz.nazwaGracza, ruch)
	return odp.err
}

func (g *gra) StanGry(graczID string) (string, error) {
	gracz, ok := g.stolik[graczID]
	if !ok {
		return "", fmt.Errorf("StanGry: nie ma gracza: %q", graczID)
	}
	log.Printf("StanGry: gra %s, gracz %q: rząda status\n", g.graID, gracz.nazwaGracza)
	stan, ok := <-gracz.kanStatus
	if !ok {
		return "", fmt.Errorf("gracz status !ok")
	}
	log.Printf("StanGry: gra %s, gracz %q: wysyłam status: %q\n", g.graID, gracz.nazwaGracza, stan)
	return stan, nil
}

func (g *gra) przebiegRozgrywki() {
	log.Printf("Rozgrywka1: gra %s: rozpoczęcie rozgrywki\n", g.graID)
	// dołączanie graczy
	timeout := time.After(DOLACZANIE_GRACZY_TIMEOUT)
	for i := 1; i <= g.liczbaGraczy; i++ {
		select {
		case req := <-g.kanDolaczGracza:
			log.Printf("gra %s: dolaczanie gracza: %s\n", g.graID, req.nazwaGracza)
			graczID := g.getNowyGraczID()
			nowyGracz := &gracz{
				graczID:     graczID,
				nazwaGracza: req.nazwaGracza,
				kanRuch:     make(chan reqRuchGracza),
				kanStatus:   make(chan string),
			}
			g.stolik[graczID] = nowyGracz
			odp := odpDolaczGracza{
				graczID: graczID,
				err:     nil,
			}
			req.kanOdp <- odp

		case <-timeout:
			g.koniec(fmt.Errorf("nie zebrał się komplet graczy w odpowiednim czasie: %v", DOLACZANIE_GRACZY_TIMEOUT))
			return
		}
	}

	// kolejność graczy alfabetycznie według losowo wygenerowanych ID
	kolejnoscGraczy := []string{}
	for graczID := range g.stolik {
		kolejnoscGraczy = append(kolejnoscGraczy, graczID)
	}
	sort.Strings(kolejnoscGraczy)
	log.Printf("Rozgrywka2: gra %s: kolejność graczy:", g.graID)
	for _, graczID := range kolejnoscGraczy {
		log.Print(" ", g.stolik[graczID].nazwaGracza)
	}

	// wykonywanie ruchów
	i := 0

	// wyślij pierwszy status
	ruszajacyGracz := g.stolik[kolejnoscGraczy[i]]
	stan := g.logika.StanGry(i)
	timeout2 := time.After(WYSLIJ_STATUS_TIMEOUT)
	select {
	case ruszajacyGracz.kanStatus <- stan:
	case <-timeout2:
		g.koniec(fmt.Errorf("upłynął czas dla gracza: %s", ruszajacyGracz.nazwaGracza))
		return
	}
	log.Printf("Rozgrywka3: gra %s: wysłano status dla gracza %q\n", g.graID, ruszajacyGracz.nazwaGracza)

	for {
		ruszajacyGracz = g.stolik[kolejnoscGraczy[i]]

		log.Printf("Rozgrywka4: gra %s: ruszający gracz %d %q\n", g.graID, i, ruszajacyGracz.nazwaGracza)

		timeout1 := time.After(RUCH_GRACZA_TIMEOUT)
		select {
		case req := <-ruszajacyGracz.kanRuch:

			err := g.logika.Ruch(i, req.ruch)
			if err != nil {
				req.kanOdp <- odpRuchGracza{err: err}
				continue
			} else {
				req.kanOdp <- odpRuchGracza{}
			}

		case <-timeout1:
			g.koniec(fmt.Errorf("upłynął czas ruchu dla gracza: %s", ruszajacyGracz.nazwaGracza))
			return
		}

		log.Printf("Rozgrywka5: gra %s: wykonano ruch gracza %q\n", g.graID, ruszajacyGracz.nazwaGracza)

		i = g.nastepny(i)
		nastepnyGracz := g.stolik[kolejnoscGraczy[i]]

		stan := g.logika.StanGry(i)

		timeout2 := time.After(WYSLIJ_STATUS_TIMEOUT)
		select {
		case nastepnyGracz.kanStatus <- stan:
		case <-timeout2:
			g.koniec(fmt.Errorf("upłynął czas dla gracza: %s", nastepnyGracz.nazwaGracza))
			return
		}
		log.Printf("Rozgrywka6: gra %s: wysłano status dla gracza %d %q\n", g.graID, i, nastepnyGracz.nazwaGracza)
	} //for
}

func (g *gra) koniec(err error) {
	log.Printf("koniec: gra %s: koniec rozgrywki: %v\n", g.graID, err)
	for _, gracz := range g.stolik {
		close(gracz.kanRuch)
	}
	close(g.kanDolaczGracza)
	g.kanKoniecGry <- err
}

func (g *gra) nastepny(i int) int {
	nast := i + 1
	if nast == g.liczbaGraczy {
		return 0
	} else {
		return nast
	}
}

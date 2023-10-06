package silnik

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/slaraz/turniej/gra_go/proto"
)

const (
	MIN_LICZBA_GRACZY         = 1
	MAX_LICZBA_GRACZY         = 5
	DOLACZANIE_GRACZY_TIMEOUT = time.Second * 300
	RUCH_GRACZA_TIMEOUT       = time.Second * 100
	WYSLIJ_STATUS_TIMEOUT     = time.Second * 100
)

type gra struct {
	graID        string
	liczbaGraczy int
	graczeByID   map[string]*gracz
	logika       ILogikaGry

	kanDolaczGracza chan reqDolaczGracza
	kanKoniecGry    chan error
}

type gracz struct {
	graczID     string
	nazwaGracza string
	kolorZolwia proto.KolorZolwia
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
		graID:           graID,
		liczbaGraczy:    liczbaGraczy,
		graczeByID:      map[string]*gracz{},
		logika:          getLogikaGry(liczbaGraczy),
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
	karta  proto.Karta
	kanOdp chan odpRuchGracza
}

type odpRuchGracza struct {
	err error
}

func (g *gra) WykonajRuch(graczID string, zagranaKarta proto.Karta) (string, error) {
	gracz, ok := g.graczeByID[graczID]
	if !ok {
		return "", fmt.Errorf("WykonajRuch: nie ma gracza: %q", graczID)
	}
	kanOdp := make(chan odpRuchGracza)
	log.Printf("WykonajRuch: gra %s, gracz %q: rząda wykonania ruchu %q\n", g.graID, gracz.nazwaGracza, zagranaKarta)
	gracz.kanRuch <- reqRuchGracza{
		karta:  zagranaKarta,
		kanOdp: kanOdp,
	}
	odp := <-kanOdp
	log.Printf("WykonajRuch: gra %s, gracz %q: wykonano ruchu %q\n", g.graID, gracz.nazwaGracza, zagranaKarta)
	return gracz.graczID, odp.err
}

func (g *gra) StanGry(graczID string) (*proto.StanGry, error) {
	gracz, ok := g.graczeByID[graczID]
	if !ok {
		return nil, fmt.Errorf("StanGry: nie ma gracza: %q", graczID)
	}
	log.Printf("StanGry: gra %s, gracz %q: rząda status\n", g.graID, gracz.nazwaGracza)
	stan, ok := <-gracz.kanStan
	// TODO: albo kanał koniec gry
	if !ok {
		return nil, fmt.Errorf("gracz status !ok")
	}
	stan.GraID = g.graID
	stan.GraczID = graczID
	stan.TwojKolor = gracz.kolorZolwia

	log.Printf("StanGry: gra %s, gracz %q: wysyłam status: %q\n", g.graID, gracz.nazwaGracza, stan)
	return stan, nil
}

func (g *gra) przebiegRozgrywki() {
	log.Printf("Rozgrywka1: gra %s: rozpoczęcie rozgrywki\n", g.graID)
	// dołączanie graczy
	gracze := []*gracz{}
	timeout := time.After(DOLACZANIE_GRACZY_TIMEOUT)
	for i := 1; i <= g.liczbaGraczy; i++ {
		select {
		case req := <-g.kanDolaczGracza:
			log.Printf("gra %s: dolaczanie gracza: %s\n", g.graID, req.nazwaGracza)
			graczID := g.getNowyGraczID()
			nowyGracz := &gracz{
				graczID:     graczID,
				kolorZolwia: g.losujKolorDlaGracza(),
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
			g.koniec(fmt.Errorf("nie zebrał się komplet graczy w odpowiednim czasie: %v", DOLACZANIE_GRACZY_TIMEOUT))
			return
		}
	}

	// losujemy kolejność graczy
	rand.Shuffle(len(gracze), func(i, j int) {
		gracze[i], gracze[j] = gracze[j], gracze[i]
	})
	log.Printf("Rozgrywka2: gra %s: kolejność graczy:", g.graID)
	for _, gracz := range gracze {
		log.Print(" ", gracz.nazwaGracza)
	}

	// wykonywanie ruchów
	i := 0

	// wyślij pierwszy status
	ruszajacyGracz := gracze[i]
	stan, _ := g.logika.GetGameStatus(i + 1)
	timeout2 := time.After(WYSLIJ_STATUS_TIMEOUT)
	select {
	case ruszajacyGracz.kanStan <- stan:
	case <-timeout2:
		g.koniec(fmt.Errorf("upłynął czas dla gracza: %s", ruszajacyGracz.nazwaGracza))
		return
	}
	log.Printf("Rozgrywka3: gra %s: wysłano status dla gracza %q\n", g.graID, ruszajacyGracz.nazwaGracza)

	for {
		ruszajacyGracz = gracze[i]

		log.Printf("Rozgrywka4: gra %s: ruszający gracz %d %q\n", g.graID, i, ruszajacyGracz.nazwaGracza)

		timeout1 := time.After(RUCH_GRACZA_TIMEOUT)
		select {
		case req := <-ruszajacyGracz.kanRuch:

			err := g.logika.Move(ruszajacyGracz.kolorZolwia, req.karta)
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
		nastepnyGracz := gracze[i]

		stan, _ := g.logika.GetGameStatus(i)

		timeout2 := time.After(WYSLIJ_STATUS_TIMEOUT)
		select {
		case nastepnyGracz.kanStan <- stan:
		case <-timeout2:
			g.koniec(fmt.Errorf("upłynął czas dla gracza: %s", nastepnyGracz.nazwaGracza))
			return
		}
		log.Printf("Rozgrywka6: gra %s: wysłano status dla gracza %d %q\n", g.graID, i, nastepnyGracz.nazwaGracza)
	} //for
}

func (g *gra) koniec(err error) {
	log.Printf("koniec: gra %s: koniec rozgrywki: %v\n", g.graID, err)
	for _, gracz := range g.graczeByID {
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

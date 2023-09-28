package silnik

import (
	"fmt"
	"sort"
	"time"

	"github.com/slaraz/turniej/gra_go/serwer/logika"
)

const (
	MIN_LICZBA_GRACZY         = 1
	MAX_LICZBA_GRACZY         = 5
	DOLACZANIE_GRACZY_TIMEOUT = time.Second * 30
	RUCH_GRACZA_TIMEOUT       = time.Second * 10
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
	// ograniczamy tworzenie gry z liczbą graczy ponad limit ustalony na serwerze
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
	kanOdp := make(chan odpRuchGracza)
	g.stolik[graczID].kanRuch <- reqRuchGracza{
		ruch:   ruch,
		kanOdp: kanOdp,
	}
	odp := <-kanOdp
	return odp.err
}

func (g *gra) StanGry(graczID string) (string, error) {
	gracz := g.stolik[graczID]
	stan, ok := <-gracz.kanStatus
	if !ok {
		return "", fmt.Errorf("gracz status !ok")
	}
	return stan, nil
}

func (g *gra) przebiegRozgrywki() {
	// dołączanie graczy
	timeout := time.After(DOLACZANIE_GRACZY_TIMEOUT)
	for i := 1; i <= g.liczbaGraczy; i++ {
		select {
		case req := <-g.kanDolaczGracza:
			graczID := g.getNowyGraczID()
			nowyGracz := &gracz{
				graczID:     graczID,
				nazwaGracza: req.nazwaGracza,
				kanRuch:     make(chan reqRuchGracza),
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

	// kolejność graczy
	kolejnoscGraczy := []string{}
	for graczID := range g.stolik {
		kolejnoscGraczy = append(kolejnoscGraczy, graczID)
	}
	// alfabetycznie według losowo wygenerowanych ID
	sort.Strings(kolejnoscGraczy)

	// wykonywanie ruchów
	i := 0
	for {

		ruszajacyGracz := g.stolik[kolejnoscGraczy[i]]

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
			g.koniec(fmt.Errorf("upłynął czas dla gracza: %s", ruszajacyGracz.graczID))
			return
		}

		i := g.nastepny(i)
		nastepnyGracz := g.stolik[kolejnoscGraczy[i]]

		stan := g.logika.StanGry(i)

		timeout2 := time.After(time.Second)
		select {
		case nastepnyGracz.kanStatus <- stan:
		case <-timeout2:
			g.koniec(fmt.Errorf("upłynął czas dla gracza: %s", ruszajacyGracz.graczID))
			return
		}
	} //for
}

func (g *gra) koniec(err error) {
	for _, gracz := range g.stolik {
		close(gracz.kanRuch)
	}
	close(g.kanDolaczGracza)
	g.kanKoniecGry <- err
	return
}

func (g *gra) nastepny(i int) int {
	nast := i + 1
	if nast == g.liczbaGraczy {
		return 0
	} else {
		return nast
	}
}

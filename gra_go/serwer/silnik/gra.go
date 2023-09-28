package silnik

import (
	"fmt"
	"sort"
	"time"

	"github.com/slaraz/turniej/gra_go/proto"
)

type gra struct {
	liczbaGraczy int
	stolik       map[string]*gracz
	gracze       []*gracz
	kanWezGracza chan string
	kanDone      chan error
}

type gracz struct {
	graczId string
	kanRuch chan reqRuch
}

const (
	MAX_LICZBA_GRACZY        = 5
	ZBIERANIE_GRACZY_TIMEOUT = time.Second * 30
	RUCH_GRACZA_TIMEOUT      = time.Second * 10
)

func nowaGra(liczbaGraczy int) (*gra, error) {
	if liczbaGraczy > MAX_LICZBA_GRACZY {
		return nil, fmt.Errorf("zbyt duża liczba graczy: %d, maksymalna dozwolona liczba graczy to: %d", liczbaGraczy, MAX_LICZBA_GRACZY)
	}

	g := &gra{
		liczbaGraczy: liczbaGraczy,
		stolik:       map[string]*gracz{},
		kanDone:      make(chan error),
	}

	go g.przebiegRozgrywki()

	return g, nil
}

type reqDodajGracza struct {
	wizytowka *proto.WizytowkaGracza
	kanOdp    chan odpDodajGracza
}

type odpDodajGracza struct {
	graczID string
	err     error
}

func (g *gra) WezGracza(wizytowka *proto.WizytowkaGracza) (string, error) {
	graczID, ok := <-g.kanWezGracza
	if !ok {
		return "", fmt.Errorf("błąd wybierania gracza")
	}
	return graczID, nil
}

type reqRuch struct {
	ruch   string
	kanOdp chan odpRuch
}

type odpRuch struct {
	err error
}

func (g *gra) przebiegRozgrywki() {

	// przygotowanie krzesełek
	for i := 1; i <= g.liczbaGraczy; i++ {
		graczID := g.getNowyGraczID()
		nowyGracz := &gracz{
			graczId: graczID,
			kanRuch: make(chan reqRuch),
		}
		g.gracze = append(g.gracze, nowyGracz)
	}

	// pobieranie graczy
	timeout := time.After(ZBIERANIE_GRACZY_TIMEOUT)
	for _, gracz := range g.gracze {
		select {
		case g.kanWezGracza <- gracz.graczId:
		case <-timeout:
			close(g.kanWezGracza)
			g.kanDone <- fmt.Errorf("nie zebrał się komplet graczy w odpowiednim czasie: %v", ZBIERANIE_GRACZY_TIMEOUT)
			return
		}
	}

	// posortowanie graczy po graczID
	sort.Slice(g.gracze, func(i, j int) bool {
		return g.gracze[i].graczId < g.gracze[j].graczId
	})

	// wykonywanie ruchów
	for {
		for _, gracz := range g.gracze {
			timeout := time.After(RUCH_GRACZA_TIMEOUT)
			select {
			case req := <-gracz.kanRuch:
				fmt.Println(req.ruch)
				req.kanOdp <- odpRuch{}
			case <-timeout:
				koniec(fmt.Errorf("upłynął czas dla gracza: %s", gracz.graczId))
				return
			}
		}
	}
}

func (g *gra) getKrzeselko() string {
	graczId := ""
	for {
		graczId = generujLosoweId(DLUGOSC_GRACZ_ID)
		// czy jest takie id?
		if _, ok := g.stolik[graczId]; !ok {
			// nie ma, bierzemy
			break
		}
	}
	return graczId
}

// func (g *gra) stanGry(graczId string) (*proto.StanGry, error) {
// 	stanGry := &proto.StanGry{
// 		GraId:             g.graId,
// 		GraczId:           graczId,
// 		SytuacjaNaPlanszy: "|__|__|",
// 		TwojeKarty:        "A1,Z7",
// 	}

// 	return stanGry, nil
// }

// func (g *gra) ruchGracza(ruch *proto.RuchGracza) error {
// 	gracz, ok := g.stolik[ruch.GraczId]
// 	if !ok {
// 		return fmt.Errorf("brak gracza %q", ruch.GraczId)
// 	}

// 	time.Sleep(time.Second)
// 	log.Printf("gracz %q zagrał karę %q\n", gracz.wizytowka.Nazwa, ruch.ZagranaKarta)
// 	return nil
// }

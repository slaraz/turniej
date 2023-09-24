package silnik

import (
	"fmt"
	"log"
	"time"

	"github.com/slaraz/turniej/gra_go/proto"
)

type gra struct {
	graId        string
	liczbaGraczy int
	stolik       map[string]*gracz
	kan          chan string
}

type gracz struct {
	graczId   string
	wizytowka *proto.WizytowkaGracza
}

func nowaGra(graId string, liczbaGraczy int) *gra {

	ticker := time.NewTicker(time.Second)
	kan := make(chan string)

	go func() {
		for {
			select {
			case s := <-kan:
				fmt.Println(s)
			case t := <-ticker.C:
				fmt.Println("tik", t)
			}
		}
	}()

	return &gra{
		graId:        graId,
		liczbaGraczy: liczbaGraczy,
		stolik:       map[string]*gracz{},
		kan:          kan,
	}
}

func (g *gra) dodajGracza(wizytowka *proto.WizytowkaGracza) (string, error) {
	if g.liczbaGraczy == len(g.stolik) {
		return "", fmt.Errorf("wszystkie miejsca zajęte, liczba miejsc: %d", len(g.stolik))
	}

	graczId := ""
	for {
		graczId = generujLosoweId(DLUGOSC_GRACZ_ID)
		// czy jest takie id?
		if _, ok := g.stolik[graczId]; !ok {
			// nie ma, bierzemy
			break
		}
	}
	nowyGracz := &gracz{
		graczId:   graczId,
		wizytowka: wizytowka,
	}
	g.stolik[graczId] = nowyGracz
	return graczId, nil
}

func (g *gra) stanGry(graczId string) (*proto.StanGry, error) {
	stanGry := &proto.StanGry{
		GraId:             g.graId,
		GraczId:           graczId,
		SytuacjaNaPlanszy: "|__|__|",
		TwojeKarty:        "A1,Z7",
	}

	return stanGry, nil
}

func (g *gra) ruchGracza(ruch *proto.RuchGracza) error {
	gracz, ok := g.stolik[ruch.GraczId]
	if !ok {
		return fmt.Errorf("brak gracza %q", ruch.GraczId)
	}

	time.Sleep(time.Second)
	log.Printf("gracz %q zagrał karę %q\n", gracz.wizytowka.Nazwa, ruch.ZagranaKarta)
	return nil
}

package silnik

import (
	"fmt"
	"log"
	"time"

	"github.com/slaraz/turniej/gra_go/proto"
)

type gra struct {
	liczbaGraczy int
	stolik       map[string]*gracz
	kanGracze    chan reqDodajGracza
}

type gracz struct {
	graczId   string
	wizytowka *proto.WizytowkaGracza
}

func nowaGra(graId string, liczbaGraczy int) *gra {

	g := &gra{
		liczbaGraczy: liczbaGraczy,
		stolik:       map[string]*gracz{},
	}

	go g.przebiegRozgrywki()

	return g
}

type reqDodajGracza struct {
	wizytowka *proto.WizytowkaGracza
	kanOdp chan odpDodajGracza
}

type odpDodajGracza struct {
	graczID string
	err error
}

func (g *gra) DodajGracza(wizytowka *proto.WizytowkaGracza) (string, error) {
	kanOdp :=make(chan odpDodajGracza)
	g.kanGracze <- reqDodajGracza{wizytowka, kanOdp}
	odp := <- kanOdp
	return odp.graczId, odp.err
}

func (g *gra) przebiegRozgrywki() err {

	// zbieranie graczy
	czasOut := time.After(time.Second*30)
	for i:=1; i<= g.liczbaGraczy; i++ {


		select {

		case req := <- g.kanGracze:
			graczId := g.dodajGracza()

			fmt.Println("dodałem gracza:", graczId)

			req.kanOdp <- odp

		case t := <-czasMinal:
			fmt.Println("czasOut:", t)
			return fmt.Errorf("Minął czas na zbieranie graczy.")
		}
	}
	// wykonywanie ruchów
	for {
		select {
		case s := <-kan:
			fmt.Println(s)
		case t := <-ticker.C:
			fmt.Println("tik", t)
		}
	}

	// wynik zakończonej gry

	// usuń grę z Areny

}

type reqDodajGracza struct {
	nazwaGracza string
	kanOdp chan odpDodajGracza
}

type odpDodajGracza struct {
	graczId string
	err error
}

func (g* gra) getKrzeselko() (string)
{
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

func (g *gra) dodajGracza(wizytowka *proto.WizytowkaGracza) string {
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
	return graczId
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

package silnik

import (
	"log"
	"math/rand"

	"github.com/slaraz/turniej/gra_go/proto"
)

// const dozwoloneZnaki = "abcdefghijklmnopqrstuvwxyz0123456789"
const (
	DLUGOSC_GRA_ID   = 2
	DLUGOSC_GRACZ_ID = 10
	DOZWOLONE_ZNAKI  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generujLosoweId(dlugosc int) string {
	id := make([]byte, dlugosc)
	for i := range id {
		id[i] = DOZWOLONE_ZNAKI[rand.Intn(len(DOZWOLONE_ZNAKI))]
	}
	return string(id)
}

func (sil *Arena) getNowaGraID() string {
	graID := ""
	for {
		graID = generujLosoweId(DLUGOSC_GRA_ID)
		// czy jest takie id?
		if _, ok := sil.aktywneGry[graID]; !ok {
			// nie ma, bierzemy
			return graID
		}
	}

}

func (g *gra) getNowyGraczID() string {
	graczID := ""
	for {
		graczID = generujLosoweId(DLUGOSC_GRACZ_ID)
		// czy jest takie id?
		if _, ok := g.graczeByID[graczID]; !ok {
			// nie ma, bierzemy
			return graczID
		}
	}
}

func (g *gra) losujKolorDlaGracza() proto.KolorZolwia {
	wykorzystaneKolory := map[proto.KolorZolwia]bool{}
	for _, gracz := range g.graczeByID {
		wykorzystaneKolory[gracz.kolorZolwia] = true
	}
	wolne := []proto.KolorZolwia{}
	for v := range proto.KolorZolwia_name {
		kolor := proto.KolorZolwia(v)
		if _, ok := wykorzystaneKolory[kolor]; !ok {
			wolne = append(wolne, kolor)
		}
	}
	if len(wolne) == 0 {
		log.Fatal("nie ma wolnych kolorów")
	}
	return wolne[rand.Intn(len(wolne))]
}

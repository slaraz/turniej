package silnik

import (
	"math/rand"
)

const dozwoloneZnaki = "abcdefghijklmnopqrstuvwxyz0123456789"

func generujLosoweId(dlugosc int) string {
	id := make([]byte, dlugosc)
	for i := range id {
		id[i] = dozwoloneZnaki[rand.Intn(len(dozwoloneZnaki))]
	}
	return string(id)
}

const (
	DLUGOSC_GRA_ID   = 4
	DLUGOSC_GRACZ_ID = 10
)

func (sil *ArenaGry) getNowaGraID() string {
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
		if _, ok := g.stolik[graczID]; !ok {
			// nie ma, bierzemy
			return graczID
		}
	}
}

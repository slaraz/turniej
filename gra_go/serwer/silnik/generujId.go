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
	DLUGOSC_GRA_ID   = 3
	DLUGOSC_GRACZ_ID = 10
)

func (sil *ArenaGry) getNowaGraID() string {
	graId := ""
	for {
		graId = generujLosoweId(DLUGOSC_GRA_ID)
		// czy jest takie id?
		if _, ok := sil.aktywneGry[graId]; !ok {
			// nie ma, bierzemy
			return graId
		}
	}

}

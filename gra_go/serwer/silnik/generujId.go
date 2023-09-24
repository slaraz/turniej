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

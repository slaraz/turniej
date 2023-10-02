package logika

import (
	"fmt"
	"log"
)

type LogikaGry struct{}

func (*LogikaGry) Ruch(nrGracza int, ruch string) error {
	log.Printf("LogikaGry: gracz nr %d zrobił ruch: %v\n", nrGracza, ruch)
	return nil
}

func (*LogikaGry) StanGry(nrGracza int) string {
	return fmt.Sprintf("stan gry dla gracza %d\n", nrGracza)
}

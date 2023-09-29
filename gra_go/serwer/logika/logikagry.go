package logika

import "fmt"

type LogikaGry struct{}

func (*LogikaGry) Ruch(nrGracza int, ruch string) error {
	fmt.Printf("gracz nr %d zrobił ruch: %v", nrGracza, ruch)
	return nil
}

func (*LogikaGry) StanGry(nrGracza int) string {
	return fmt.Sprintf("stan gry dla gracza %d", nrGracza)
}
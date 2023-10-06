package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/slaraz/turniej/gra_go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	IP_ADDR               = "localhost:50051"
	NOWY_MECZ_TIMEOUT     = time.Second * 5
	DOLACZ_DO_GRY_TIMEOUT = time.Second * 150
	RUCH_GRACZA_TIMEOUT   = time.Second * 150
)

var (
	addr  = flag.String("addr", IP_ADDR, "adres serwera gry")
	nazwa = flag.String("nazwa", "Ziutek", "nazwa gracza")
	nowa  = flag.Bool("nowa", false, "tworzy nową grę na serwerze")
	graID = flag.String("gra", "", "dołącza do gry o podanym id")
)

func main() {
	fmt.Println("Start")
	defer fmt.Println("Koniec.")

	flag.Parse()

	// Utowrzenie połączenia z serwerem gry.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial: %v", err)
	}
	defer conn.Close()
	c := proto.NewGraClient(conn)

	conn.GetState()
	// Jeśli podano opcję -nowa, to utwórz nową grę.
	if *nowa {
		ctx, cancel := context.WithTimeout(context.Background(), NOWY_MECZ_TIMEOUT)
		defer cancel()
		nowaGraInfo, err := c.NowyMecz(ctx, &proto.KonfiguracjaGry{LiczbaGraczy: 2})
		if err != nil {
			log.Fatalf("c.NowyMecz: %v", err)
		}
		log.Printf("Nowa gra %q\n", nowaGraInfo.GraID)

		*graID = nowaGraInfo.GraID
	}

	// Jeśli nie utworzono -nowa,
	// ani nie podano opcji -gra, to kończymy.
	if *graID == "" {
		flag.Usage()
		return
	}

	// Dołączamy do gry graID.
	stanGry := dolaczDoGry(c, *graID, *nazwa)
	fmt.Printf("Stan gry: plansza: %v, karty: %v\n", stanGry.Plansza, stanGry.TwojeKarty)

	// Przebieg gry.
	for {

		// Gracz podaje kartę na konsoli.
		fmt.Println("Wybierz kartę do zagrania:")
		karta := wczytajKarte()

		// Wysyłam ruch do serwera.
		stanGry = wyslijRuch(c, &proto.RuchGracza{
			GraID:        stanGry.GraID,
			GraczID:      stanGry.GraczID,
			ZagranaKarta: karta,
		})
		fmt.Printf("Stan gry: plansza: %v, karty: %v\n", stanGry.Plansza, stanGry.TwojeKarty)

		if stanGry.CzyKoniec {
			fmt.Println("Koniec gry, wygrał gracz nr", stanGry.KtoWygral)
			break
		}
	}
}

func wczytajKarte() proto.Karta {
	fmt.Print("> ")
	var karta string
	_, err := fmt.Scanln(&karta)
	if err != nil {
		log.Fatalf("Błąd wczytywania karty: %v", err)
	}
	k, ok := proto.Karta_value[karta]
	if !ok {
		log.Fatalf("Niepoprawna karta: %q", karta)
	}
	return proto.Karta(k)
}

func dolaczDoGry(c proto.GraClient, graID, nazwa string) *proto.StanGry {
	log.Printf("Gracz %s dołącza do gry %q", nazwa, graID)
	ctx, cancel := context.WithTimeout(context.Background(), DOLACZ_DO_GRY_TIMEOUT)
	defer cancel()
	log.Println("Czekam na stan gry...")
	stanGry, err := c.DolaczDoGry(ctx, &proto.Dolaczanie{
		GraID:       graID,
		NazwaGracza: nazwa,
	})
	if err != nil {
		log.Fatalf("c.Dolacz: %v", err)
	}
	return stanGry
}

func wyslijRuch(c proto.GraClient, ruch *proto.RuchGracza) *proto.StanGry {
	log.Printf("Gracz %s-%s zagrywa kartę: %v", ruch.GraID, ruch.GraczID, ruch.ZagranaKarta)
	ctx, cancel := context.WithTimeout(context.Background(), RUCH_GRACZA_TIMEOUT)
	defer cancel()
	log.Println("Czekam na stan gry...")

	stanGry, err := c.MojRuch(ctx, ruch)
	if err != nil {
		log.Fatalf("c.MojRuch: %v", err)
	}
	return stanGry
}

func koniecGry(stanGry *proto.StanGry) bool {
	return stanGry.CzyKoniec
}

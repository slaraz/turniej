package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/slaraz/turniej/gra_go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	IP_ADDR               = "localhost:50051"
	NOWY_MECZ_TIMEOUT     = time.Second * 5
	DOLACZ_DO_GRY_TIMEOUT = time.Second * 1000
	RUCH_GRACZA_TIMEOUT   = time.Second * 1000
)

var (
	addr  = flag.String("addr", IP_ADDR, "adres serwera gry")
	nazwa = flag.String("nazwa", "Ziutek", "nazwa gracza")
	nowa  = flag.Bool("nowa", false, "tworzy nową grę na serwerze")
	graID = flag.String("gra", "", "dołącza do gry o podanym id")
	lg    = flag.Int("lg", 2, "określa liczbę graczy")
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
		nowaGraInfo, err := c.NowyMecz(ctx, &proto.KonfiguracjaGry{LiczbaGraczy: int32(*lg)})
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

	var (
		kartyDlaKtorychTrzebaPodacKolor = map[proto.Karta]bool{
			proto.Karta_L1:  true,
			proto.Karta_L2:  true,
			proto.Karta_A1:  true,
			proto.Karta_A1B: true,
		}
		karta proto.Karta
		kolor proto.KolorZolwia
	)

	// przebieg gry

	// dołączamy do gry graID
	stanGry := dolaczDoGry(c, *graID, *nazwa)
	for {
		// wypisuję stan gry na ekranie
		drukujStatus(stanGry)
		if stanGry.CzyKoniec {
			return
		}
		for {
			// gracz podaje kartę na konsoli
			karta = wczytajKarte()

			if _, ok := kartyDlaKtorychTrzebaPodacKolor[karta]; ok {
				kolor = wczytajKolor()
			} else {
				kolor = proto.KolorZolwia_XXX
			}

			// wysyłam ruch do serwera
			nowyStan, err := wyslijRuch(c, &proto.RuchGracza{
				GraID:        stanGry.GraID,
				GraczID:      stanGry.GraczID,
				ZagranaKarta: karta,
				KolorWybrany: kolor,
			})
			if err != nil && status.Code(err) == codes.InvalidArgument {
				// zły ruch
				fmt.Printf("Błąd ruchu: %v\n", err)
				continue
			} else if err != nil {
				// inny błąd, np. połączenie z serwerem
				log.Fatalf("wyslijRuch: status: %v, err: %v", status.Code(err), err)
			}
			// ruch ok
			stanGry = nowyStan
			break
		}
	}
}

func wczytajKolor() proto.KolorZolwia {
	fmt.Print("Wybierz kolor\n> ")
	var kolor string
	_, err := fmt.Scanln(&kolor)
	if err != nil {
		log.Fatalf("Błąd wczytywania koloru: %v", err)
	}
	k, ok := proto.KolorZolwia_value[strings.ToUpper(kolor)]
	if !ok {
		log.Fatalf("Niepoprawny kolor: %q", kolor)
	}
	return proto.KolorZolwia(k)
}

func wczytajKarte() proto.Karta {
	var karta proto.Karta
	for {
		fmt.Print("Wybierz kartę do zagrania:\n> ")
		var kartatxt string
		_, err := fmt.Scanln(&kartatxt)
		if err != nil {
			log.Fatalf("Błąd wczytywania karty: %v", err)
		}
		k, ok := proto.Karta_value[strings.ToUpper(kartatxt)]
		if !ok {
			fmt.Printf("Niepoprawna karta: %q\n", kartatxt)
			continue
		}
		karta = proto.Karta(k)
		break
	}
	return karta
}

func dolaczDoGry(c proto.GraClient, graID, nazwa string) *proto.StanGry {
	log.Printf("Gracz %s dołącza do gry %q", nazwa, graID)
	ctx, cancel := context.WithTimeout(context.Background(), DOLACZ_DO_GRY_TIMEOUT)
	defer cancel()
	log.Println("Czekam na odpowiedź od serwera...")
	stanGry, err := c.DolaczDoGry(ctx, &proto.Dolaczanie{
		GraID:       graID,
		NazwaGracza: nazwa,
	})
	if err != nil {
		log.Fatalf("c.Dolacz: %v", err)
	}
	return stanGry
}

func wyslijRuch(c proto.GraClient, ruch *proto.RuchGracza) (*proto.StanGry, error) {
	log.Printf("Gracz %s-%s zagrywa kartę: %v", ruch.GraID, ruch.GraczID, ruch.ZagranaKarta)
	ctx, cancel := context.WithTimeout(context.Background(), RUCH_GRACZA_TIMEOUT)
	defer cancel()
	log.Println("Czekam na odpowiedź od serwera (ruch przeciwnika)...")

	return c.MojRuch(ctx, ruch)
}

func drukujStatus(stanGry *proto.StanGry) {
	if stanGry.CzyKoniec {
		fmt.Println("Koniec gry, wygrał gracz nr", stanGry.KtoWygral)
	} else {
		fmt.Printf("Twój kolor: %v, Pola:", stanGry.TwojKolor)
		for _, pole := range stanGry.Plansza {
			fmt.Printf(" %v", pole.Zolwie)
		}
		fmt.Printf(", Twoje karty: %v\n", stanGry.TwojeKarty)
	}
}

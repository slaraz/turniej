package main

import (
	"context"
	"flag"
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
	log.Println("Start")
	defer log.Println("Koniec.")

	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGraClient(conn)

	if *nowa {
		ctx, cancel := context.WithTimeout(context.Background(), NOWY_MECZ_TIMEOUT)
		defer cancel()
		nowaGraInfo, err := c.NowyMecz(ctx, &proto.KonfiguracjaGry{LiczbaGraczy: 2})
		if err != nil {
			log.Fatalf("c.NowyMecz: %v", err)
		}
		log.Printf("Nowa gra: %q\n", nowaGraInfo.GraID)

		*graID = nowaGraInfo.GraID
	}

	if *graID == "" {
		flag.Usage()
		return
	}

	log.Printf("Dołączam do gry %q\n", *graID)
	ctx, cancel := context.WithTimeout(context.Background(), DOLACZ_DO_GRY_TIMEOUT)
	defer cancel()
	dol := &proto.Dolaczanie{
		GraID: *graID,
		Wizytowka: &proto.WizytowkaGracza{
			Nazwa: *nazwa,
		},
	}
	stanGry, err := c.DolaczDoGry(ctx, dol)
	if err != nil {
		log.Fatalf("c.Dolacz: %v", err)
	}

	log.Printf("Stan gry: Gra: %q, gracz: %q", stanGry.GraID, stanGry.GraczID)

	for {
		err := ruch(c, stanGry)
		if err != nil {
			log.Fatalf("ruch: %v\n", err)
		}
	}
}

func ruch(c proto.GraClient, stanGry *proto.StanGry) error {
	ctx, cancel := context.WithTimeout(context.Background(), RUCH_GRACZA_TIMEOUT)
	defer cancel()

	rg := &proto.RuchGracza{
		GraID:        stanGry.GraID,
		GraczID:      stanGry.GraczID,
		ZagranaKarta: "A5",
	}

	log.Printf("Wykonuję ruch: Gra: %q, gracz: %q", stanGry.GraID, stanGry.GraczID)
	stanGry, err := c.MojRuch(ctx, rg)
	if err != nil {
		return err
	}

	log.Printf("Stan gry: plansza: %q, karty: %q\n", stanGry.SytuacjaNaPlanszy, stanGry.TwojeKarty)
	return nil
}

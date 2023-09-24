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

const IP_ADDR = "localhost:50051"

var (
	addr  = flag.String("addr", IP_ADDR, "adres serwera gry")
	join  = flag.String("join", "", "id gry, do której dołączyć")
	nazwa = flag.String("nazwa", "Ziutek", "id gry, do której dołączyć")
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

	wizytowka := &proto.WizytowkaGracza{
		Nazwa: *nazwa,
	}

	log.Println(*join)

	stanGry := &proto.StanGry{}
	if *join == "" {
		log.Println("Tworzę nową grę...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		stanGry, err = c.NowyMecz(ctx, wizytowka)
		if err != nil {
			log.Fatalf("c.NowyMecz: %v", err)
		}
	} else {
		log.Printf("Dołączam do gry %q\n", *join)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		dol := &proto.Dolaczanie{
			GraId:     *join,
			Wizytowka: wizytowka,
		}
		stanGry, err = c.Dolacz(ctx, dol)
		if err != nil {
			log.Fatalf("c.Dolacz: %v", err)
		}
	}
	log.Printf("Gra: %q, gracz: %q", stanGry.GraId, stanGry.GraczId)

	for {
		err := ruch(c, stanGry)
		if err != nil {
			log.Fatalf("ruch: %v\n", err)
		}
	}
}

func ruch(c proto.GraClient, stanGry *proto.StanGry) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	rg := &proto.RuchGracza{
		GraId:        stanGry.GraId,
		GraczId:      stanGry.GraczId,
		ZagranaKarta: "A5",
	}
	stanGry, err := c.MojRuch(ctx, rg)
	if err != nil {
		return err
	}

	log.Printf("plansza: %q, karty: %q\n", stanGry.SytuacjaNaPlanszy, stanGry.TwojeKarty)
	return nil
}

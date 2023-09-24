package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/slaraz/turniej/gra_go/proto"
	"google.golang.org/grpc"
)

const IP_PORT = 50051

type serwer struct {
	proto.UnimplementedGraServer
}

func (s *serwer) NowyMecz(ctx context.Context, in *proto.WizytowkaGracza) (*proto.StanGry, error) {
	log.Printf("Otrzymałem: %v\n", in.Nazwa)
	sg := proto.StanGry{
		IdGry:             "1234",
		SytuacjaNaPlanszy: "|__|__|",
		TwojeKarty:        "A1,Z7",
	}
	return &sg, nil
}

func (s *serwer) MojRuch(ctx context.Context, in *proto.RuchGracza) (*proto.StanGry, error) {
	log.Printf("Otrzymałem: %v\n", in.ZagranaKarta)
	sg := proto.StanGry{
		IdGry:             "456",
		SytuacjaNaPlanszy: "|__|__|",
		TwojeKarty:        "A1,Z7",
	}
	return &sg, nil
}

func main() {
	log.Println("Start")
	defer log.Println("Koniec.")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", IP_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterGraServer(s, &serwer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

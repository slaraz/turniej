package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/slaraz/turniej/gra_go/proto"
	"github.com/slaraz/turniej/gra_go/serwer/silnik"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const IP_PORT = 50051

type serwer struct {
	proto.UnimplementedGraServer
	arena *silnik.ArenaGry
}

func (s *serwer) NowyMecz(ctx context.Context, wizytowka *proto.WizytowkaGracza) (*proto.StanGry, error) {
	iluGraczy := 2

	graId, err := s.arena.NowaGra(iluGraczy)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	graczId, err := s.arena.DodajGraczaDoGry(graId, wizytowka)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	stanGry, err := s.arena.StanGry(graId, graczId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return stanGry, nil
}

func (s *serwer) Dolacz(ctx context.Context, dolacz *proto.Dolaczanie) (*proto.StanGry, error) {
	graId := dolacz.GraId

	graczId, err := s.arena.DodajGraczaDoGry(graId, dolacz.Wizytowka)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	stanGry, err := s.arena.StanGry(graId, graczId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return stanGry, nil
}

func (s *serwer) MojRuch(ctx context.Context, in *proto.RuchGracza) (*proto.StanGry, error) {
	err := s.arena.RuchGracza(in)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	stanGry, err := s.arena.StanGry(in.GraId, in.GraczId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return stanGry, nil
}

func main() {
	log.Println("Start")
	defer log.Println("Koniec.")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", IP_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	serwerGry := &serwer{
		arena: silnik.NowaArena(),
	}
	proto.RegisterGraServer(s, serwerGry)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

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

	gra, err := s.arena.GetGra(graId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	graczID, err := gra.WezGracza(wizytowka)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	stanGry, err := gra.StanGry(graczID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	sg := &proto.StanGry{
		GraId: graId,
		GraczId: graczID,
		SytuacjaNaPlanszy: stanGry,
		TwojeKarty: "A5,T9",
	}
	return sg, nil
}

func (s *serwer) Dolacz(ctx context.Context, dolacz *proto.Dolaczanie) (*proto.StanGry, error) {
	gra, err := s.arena.GetGra(dolacz.GraId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	graczID, err := gra.WezGracza(dolacz.Wizytowka)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	stanGry, err := gra.StanGry(graczID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	sg := &proto.StanGry{
		GraId: dolacz.GraId,
		GraczId: graczID,
		SytuacjaNaPlanszy: stanGry,
		TwojeKarty: "X7,H5",
	}
	return sg, nil
}

func (s *serwer) MojRuch(ctx context.Context, ruch *proto.RuchGracza) (*proto.StanGry, error) {
	gra, err := s.arena.GetGra(ruch.GraId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err := gra.RuchGracza(ruch)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	stanGry, err := s.arena.StanGry(ruch.GraId, ruch.GraczId)
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

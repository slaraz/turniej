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

const (
	IP_PORT = 50051
)

type serwer struct {
	proto.UnimplementedGraServer
	arena *silnik.ArenaGry
}

func (s *serwer) NowyMecz(ctx context.Context, conf *proto.KonfiguracjaGry) (*proto.NowaGraInfo, error) {
	graID, err := s.arena.NowaGra(int(conf.LiczbaGraczy))
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	info := &proto.NowaGraInfo{
		GraID: graID,
		Opis:  fmt.Sprintf("nowa graID=%q dla %d graczy", graID, conf.LiczbaGraczy),
	}
	return info, nil
}

func (s *serwer) DolaczDoGry(ctx context.Context, dolacz *proto.Dolaczanie) (*proto.StanGry, error) {
	gra, err := s.arena.GetGra(dolacz.GraID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	graczID, err := gra.DolaczGracza(dolacz.Wizytowka.Nazwa)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	log.Printf("DolaczDoGry: gracz %q dołączył do gry %q", dolacz.Wizytowka.Nazwa, dolacz.GraID)

	stanGry, err := gra.StanGry(graczID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("gra.StanGry: %v", err.Error()))
	}

	pstanGry := &proto.StanGry{
		GraczID:           graczID,
		GraID:             dolacz.GraID,
		SytuacjaNaPlanszy: stanGry,
		TwojeKarty:        "A4,X8",
	}

	return pstanGry, nil
}

func (s *serwer) MojRuch(ctx context.Context, ruch *proto.RuchGracza) (*proto.StanGry, error) {
	gra, err := s.arena.GetGra(ruch.GraID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("arena.GetGra: %v", err.Error()))
	}

	err = gra.WykonajRuch(ruch.GraczID, ruch.ZagranaKarta)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stanGry, err := gra.StanGry(ruch.GraczID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("gra.StanGry: %v", err.Error()))
	}

	pstanGry := &proto.StanGry{
		GraID:             ruch.GraID,
		SytuacjaNaPlanszy: stanGry,
		TwojeKarty:        "A4,X8",
	}

	return pstanGry, nil
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

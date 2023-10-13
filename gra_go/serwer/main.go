package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
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
	arena *silnik.Arena
}

func (s *serwer) NowyMecz(ctx context.Context, conf *proto.KonfiguracjaGry) (*proto.NowaGraInfo, error) {
	graID, err := s.arena.NowaGra(int(conf.LiczbaGraczy), conf.GraID)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	info := &proto.NowaGraInfo{
		GraID: graID,
	}
	return info, nil
}

func (s *serwer) DolaczDoGry(ctx context.Context, dolacz *proto.Dolaczanie) (*proto.StanGry, error) {
	gra, err := s.arena.GetGra(dolacz.GraID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	graczID, err := gra.DolaczGracza(dolacz.NazwaGracza)
	if err != nil {
		return nil, status.Error(codes.ResourceExhausted, err.Error())
	}

	log.Printf("%s DolaczDoGry(): dołączył gracz %q", dolacz.GraID, dolacz.NazwaGracza)

	stanGry, err := gra.StanGry(graczID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("gra.StanGry: %v", err.Error()))
	}

	return stanGry, nil
}

func (s *serwer) MojRuch(ctx context.Context, ruch *proto.RuchGracza) (*proto.StanGry, error) {
	gra, err := s.arena.GetGra(ruch.GraID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("arena.GetGra: %v", err.Error()))
	}

	graczID, err := gra.WykonajRuch(ruch.GraczID, ruch)
	if err != nil {
		log.Printf("MojRuch: gra.WykonajRuch: %v", err.Error())
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stanGry, err := gra.StanGry(graczID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("gra.StanGry: %v", err.Error()))
	}

	return stanGry, nil
}

func main() {
	log.Println("Start")
	defer log.Println("Koniec.")

	//const ELASTIC = "tylko na konsole"
	//const ELASTIC = "url=http://test:testtest@192.168.88.100:9200; indeks=hackaton"
	// logowanie do Elastica

	cert, _ := ioutil.ReadFile("/home/rtuser/ca.crt")
	const INDEX = "hackaton"
	cfg := elasticsearch.Config{
		Addresses: []string{"https://test:testtest@localhost:9200"},
		CACert:    cert,
		Username:  "test",
		Password:  "testtest",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("[Elastic] błąd elasticsearch.NewClient(): %q", err)
	}

	es.Indices.Create(INDEX)

	document := struct {
		Name string `json:"name"`
	}{
		"go-elasticsearch",
	}
	data, _ := json.Marshal(document)
	res, err := es.Index(INDEX, bytes.NewReader(data))

	// res, err := es.Index(
	// 	INDEX,
	// 	strings.NewReader(`{"title" : "Loguje się do Elastica"}`),
	// )
	if err != nil {
		log.Printf("[Elastic] błąd es.klient.Index(): %v", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("[Elastic] res.StatusCode: %d, ioutil.ReadAll(res.Body): %v", res.StatusCode, err)
		}
		log.Printf("[Elastic] res.StatusCode: %d, resp: %v", res.StatusCode, string(body))
	}

	//kl.Indeks("hackaton")
	//kl.Loguj(map[string]interface{}{
	// logg := maplogger.NowyLogger(ELASTIC)
	// defer logg.Close()
	// // logowanie Elastic
	// mlog := logg.NowyWpis()
	// mlog["hello"] = "world"
	// logg.Loguj(mlog)

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

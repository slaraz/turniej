package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	pb "github.com/slaraz/turniej/go/gra_proto"
)

type server struct {
	pb.UnimplementedGraServer
}

func (s *server) SayHello(ctx context.Context, in *pb.WizytowkaGracza) (*pb.StanGry, error) {
	log.Printf("Received: %v", in.GetNazwa())

	idGry := uuid.New().String()
	sg := pb.StanGry{
		IdGry: idGry,
	}

	return &sg, nil
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}

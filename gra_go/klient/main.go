package main

import (
	"context"
	"log"
	"time"

	pb "github.com/slaraz/turniej/gra_go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const IP_ADDR = "localhost:50051"

func main() {
	log.Println("Start")
	defer log.Println("Koniec.")

	// Set up a connection to the server.
	conn, err := grpc.Dial(IP_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGraClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	r, err := c.NowyMecz(ctx, &pb.WizytowkaGracza{Nazwa: "Ziutek"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetIdGry())
}

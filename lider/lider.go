package main

import (
	"context"
	"log"
	"time"

	pb "example.com/go-msgs-grpc/msgs"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPozoManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetMontoAcumulado(ctx, &pb.NewMontoAcumulado{})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
	log.Printf("monto acumulado: %d", r.GetMontoActual())
}

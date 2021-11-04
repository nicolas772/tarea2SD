package main

import (
	"context"
	"log"
	"net"

	pb "example.com/go-msgs-grpc/msgs"
	"google.golang.org/grpc"
)

const (
	port = ":50010"
)

func NewPozoServer(mount int64) *PozoServer {
	return &PozoServer{
		monto: mount,
	}
}

type PozoServer struct {
	pb.UnimplementedPozoManagementServer
	monto int64
}

func (server *PozoServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPozoManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *PozoServer) GetMontoAcumulado(ctx context.Context, in *pb.NewMontoAcumulado) (*pb.MontoAcumulado, error) {
	log.Printf("Pozo recibe peticion monto acumulado")
	return &pb.MontoAcumulado{MontoActual: int32(server.monto)}, nil
}

func main() {
	monto_actual := 1000000
	var pozo_server *PozoServer = NewPozoServer(int64(monto_actual))
	if err := pozo_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net"

	pb "example.com/go-msgs-grpc/msgs"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

const (
	port = ":50010"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

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
	//if err := pozo_server.Run(); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}

	//aqui empieza rabbit

	conn, err := amqp.Dial("amqp://admin:test@10.6.40.194:5672/qa1")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	go func() {
		if err := pozo_server.Run(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-forever

}

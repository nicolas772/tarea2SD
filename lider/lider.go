package main

import (
	"context"
	"log"
	"time"

	pb "example.com/go-msgs-grpc/msgs"
	"github.com/streadway/amqp" //para rabbit
	"google.golang.org/grpc"
)

const (
	address = "10.6.40.194:50010"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

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
		log.Fatalf("No se pudo obtener el monto acumulado: %v", err)
	}
	log.Printf("monto acumulado: %d", r.GetMontoActual())

	//aqui comienza Rabbit

	conn_pozo, err := amqp.Dial("amqp://admin:test@10.6.40.194:5672/qa1")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn_pozo.Channel()
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

	body := "Jugador 1 Ronda 1 600000000"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

}

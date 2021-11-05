package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "example.com/go-msgs-grpc/msgs"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

const (
	port = ":50010"
	path = "./eliminados.txt"
)

func crearArchivo() {
	//Verifica que el archivo existe
	var _, err = os.Stat(path)
	//Crea el archivo si no existe
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if existeError(err) {
			return
		}
		defer file.Close()
	}
	fmt.Println("File Created Successfully", path)
}
func escribeArchivo(texto string) {
	// Abre archivo usando permisos READ & WRITE
	var file, err = os.OpenFile(path, os.O_RDWR, 0666)
	if existeError(err) {
		return
	}
	defer file.Close()
	// Escribe algo de texto linea por linea
	_, err = file.WriteString(texto + "\n")
	if existeError(err) {
		return
	}
	// Salva los cambios
	err = file.Sync()
	if existeError(err) {
		return
	}
	fmt.Println("Archivo actualizado existosamente.")
}
func existeError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

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
	monto_actual := 0
	crearArchivo()

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

	forever := make(chan bool) //while infinito

	//sub-rutina para recibir mensaje de eliminacion de jugador
	go func() {
		for d := range msgs {
			log.Printf("Nuevo Jugador Eliminado: %s", d.Body)
			monto_actual += 100000000
			texto := string(d.Body) + " " + strconv.Itoa(monto_actual)
			escribeArchivo(texto)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	//sub-rutina para informar al lider el monto actual acumulado
	go func() {
		var pozo_server *PozoServer = NewPozoServer(int64(monto_actual))
		if err := pozo_server.Run(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-forever

}

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
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

func printmenuInicial() string {
	menu :=
		`
		Bienvenido al juego del calamar (prueba Lider - Pozo)
		[ 1 ] Pedir Monto acumulado al pozo
		[ 2 ] Eliminar un jugador
		[ otro ] Salir
		Selecciona una opcion
	`
	fmt.Print(menu)
	reader := bufio.NewReader(os.Stdin)

	entrada, _ := reader.ReadString('\n')          // Leer hasta el separador de salto de línea
	eleccion := strings.TrimRight(entrada, "\r\n") // Remover el salto de línea de la entrada del usuario
	return eleccion
}
func pedirNombreJugador() string {
	menu := "Escriba el nombre del jugador\n"

	fmt.Print(menu)
	reader := bufio.NewReader(os.Stdin)

	entrada, _ := reader.ReadString('\n')          // Leer hasta el separador de salto de línea
	eleccion := strings.TrimRight(entrada, "\r\n") // Remover el salto de línea de la entrada del usuario
	return eleccion
}
func pedirRondaJugador() string {
	menu := "Escriba la ronda en la que fue eliminado el jugador (por ejemplo, ronda 1)\n"

	fmt.Print(menu)
	reader := bufio.NewReader(os.Stdin)

	entrada, _ := reader.ReadString('\n')          // Leer hasta el separador de salto de línea
	eleccion := strings.TrimRight(entrada, "\r\n") // Remover el salto de línea de la entrada del usuario
	return eleccion
}

func main() {
	//Conexion con Pozo Rabbitmq

	conn_pozo, err := amqp.Dial("amqp://admin:test@10.6.40.194:5672/qa1")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn_pozo.Close()

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

	//Conexion con Pozo gRPC
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPozoManagementClient(conn)

	//Programa

	flag := 1

	for flag == 1 {
		eleccion := printmenuInicial()

		if eleccion == "1" {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.GetMontoAcumulado(ctx, &pb.NewMontoAcumulado{})
			if err != nil {
				log.Fatalf("No se pudo obtener el monto acumulado: %v", err)
			}
			log.Printf("Monto acumulado: %d", r.GetMontoActual())
		}
		if eleccion == "2" {
			nombre := pedirNombreJugador()
			ronda := pedirRondaJugador()
			body := nombre + " " + ronda //este es el mensaje que se envia por rabbit
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
		if eleccion != "1" && eleccion != "2" {
			log.Printf("Hasta Pronto!")
			flag = 0
		}
	}

}

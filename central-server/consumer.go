// consumer.go

package main

import (
	"encoding/json"
	"time"
	"fmt"
	"log"
	"os"
	"github.com/streadway/amqp"
)

/* ad

 GET /Suscribe {queue: "name:5672"}
 Levantar hilo que ejecute main.
 

	
Address */

var DB_CONN_STRING = os.Getenv("DB_CONN_STRING")

type Medicion struct {
	Datetime time.Time
    Sensor string
	Sector string
    Presion int
}

// Thread
func main() {

	/* //Create DB connection
	db, err := sql.Open("postgres", DB_CONN_STRING)
	if err != nil {
		log.Fatal(err)
	} */

	conn, err := amqp.Dial("amqp://queue:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"queue", // queue name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var medicion Medicion
			err := json.Unmarshal(d.Body, &medicion)
			if err != nil {
				log.Fatalf("Failed to unmarshal medicion: %s", err)
			}

			fmt.Println("Fecha y hora: " + fmt.Sprint(medicion.Datetime) + " - Sensor: " + medicion.Sensor + " - Sector: " + medicion.Sector + " - Medicion: " + fmt.Sprint(medicion.Presion))
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
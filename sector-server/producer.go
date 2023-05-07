// producer.go

package main

import (
    "encoding/json"
    "fmt"
    "log"
	"time"
    "github.com/streadway/amqp"
)

type Medicion struct {
	Datetime time.Time
    Sensor string
	Sector string
    Presion int
}

func main() {
	time.Sleep(time.Second)
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


	for i := 1; i <= 10; i++ {
		medicion := Medicion{Datetime: time.Now(), Sensor: "sensor1", Sector: "sectorA", Presion: i}
		body, err := json.Marshal(medicion)
		if err != nil {
			log.Fatalf("Error serializing medicion to JSON: %v", err)
		}

		err = ch.Publish(
			"",            // exchange
			q.Name,        // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)

		if err != nil {
				log.Fatalf("Failed to publish a message: %v", err)
			}

		fmt.Printf("Message sent: %s", body)
	}
}
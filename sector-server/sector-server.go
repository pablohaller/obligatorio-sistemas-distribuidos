// producer.go

package main

import (
    "encoding/json"
    "fmt"
    "log"
	"time"
    "github.com/streadway/amqp"
	"net/http"
	"io"
	"strconv"
)

type Medicion struct {
	Datetime time.Time
    Sensor string
	Sector string
    Presion int
}

func main() {

	clientHealth := &http.Client{}
    req, _ := http.NewRequest("GET", "http://central-server:8080/healthcheck", nil)

    for {
        resp, err := clientHealth.Do(req)
        if err != nil {
            fmt.Println(err)
        } else if resp.StatusCode == 200 {
            fmt.Println("Status code 200 received")
            break
        }
        time.Sleep(time.Second * 5) // espera 5 segundos antes de volver a intentar
    }
	
	// suscribirse para obtener el nombre de la queue
	req, http_err := http.NewRequest("GET", "http://central-server:8080/Sector/Suscribe", nil)
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud y capturar la respuesta y el posible error
	client := &http.Client{}
	resp, http_err := client.Do(req)
	if http_err != nil {
		log.Fatalf("Error al pegar al endpoint de mediciones: %v", http_err)
	}
	defer resp.Body.Close()

	body, body_err := io.ReadAll(resp.Body)
	if body_err != nil {
		log.Fatalf("Error al leer el cuerpo de la respuesta: %v", body_err)
	}

	// Eliminar las barras diagonales
	cleanBody, err1 := strconv.Unquote(string(body))
	if err1 != nil {
		log.Fatalf("Error al limpiar el cuerpo de la respuesta: %v", err1)
	}

	// Utilizar el cuerpo de respuesta limpio
	queue := cleanBody
	var conn *amqp.Connection
	var err error

	for conn == nil {
		conn, err = amqp.Dial("amqp://" + queue + "/")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v", err)
			time.Sleep(5 * time.Second) // wait 5 seconds before retrying
		}
	}
	defer conn.Close()

// continue with other code using `conn` object

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
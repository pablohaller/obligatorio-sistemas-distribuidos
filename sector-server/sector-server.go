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
	"github.com/gin-gonic/gin"
	"os"
	"bytes"
)

type Medicion struct {
    Datetime time.Time `json:"Datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"Sensor"`
    Sector   string    `json:"Sector"`
    Presion  int       `json:"Presion"`
}

var SECTOR_NAME = os.Getenv("SECTOR_NAME")
var QUEUE_HOST = os.Getenv("QUEUE_HOST")
var channelQueue *amqp.Channel
var queue amqp.Queue

func main() {

	r := gin.Default()

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

	queueName := QUEUE_HOST
	var conn *amqp.Connection
	var err error

	for conn == nil {
		conn, err = amqp.Dial("amqp://" + queueName + ":5672/")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v", err)
			time.Sleep(5 * time.Second) // wait 5 seconds before retrying
		}
	}
	defer conn.Close()

// continue with other code using `conn` object

	channelQueue, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer channelQueue.Close()

	queue, err = channelQueue.QueueDeclare(
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
	
	go alertCentral()

	r.PUT("/Medicion", func(c *gin.Context) {
		var medicion Medicion
		// Deserializar el body JSON en la struct Medicion
		if err := c.ShouldBindJSON(&medicion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, err := json.Marshal(medicion)
		if err != nil {
			log.Fatalf("Error serializing medicion to JSON: %v", err)
		}

		err = channelQueue.Publish(
			"",            // exchange
			queue.Name,        // routing key
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

		// guardar en MONGO DB LA MEDICION


		return
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

type Suscribe struct {
    Sector   string    `json:"Sector"`
    Queue    string    `json:"Queue"`
}

func alertCentral(){
	for{

		suscribe := Suscribe{
			Queue:  QUEUE_HOST,
			Sector: SECTOR_NAME,
		}

		// Convertir la estructura a JSON
		jsonData, marshalerr := json.Marshal(suscribe)
		if marshalerr != nil {
			log.Fatal("Error al convertir la estructura a JSON:", marshalerr)
		}

		// suscribirse para obtener el nombre de la queue
		req, http_err := http.NewRequest("POST", "http://central-server:8080/Sector/Suscribe", bytes.NewBuffer(jsonData))
		if http_err != nil {
			log.Fatal("Error al crear la solicitud:", http_err)
		}
		req.Header.Set("Content-Type", "application/json")

		// AGREGAR SI ES DESEADO EL MANDAR TODOS LOS DATOS GUARDADOS DE HACE 5 MINUTOS HASTA 
		// AHORA DE MONGO DB

		time.Sleep(time.Minute * 5)
	}
}
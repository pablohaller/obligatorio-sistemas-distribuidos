// producer.go

package main

import (
    "encoding/json"
    "fmt"
    "log"
	"time"
    "github.com/streadway/amqp"
	"net/http"
	"github.com/gin-gonic/gin"
	"os"
	"bytes"
)

type Measurement struct {
    Datetime time.Time `json:"datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"sensor"`
    Sector   string    `json:"sector"`
    Pressure  float32      `json:"pressure"`
}

type Sensor struct {
    Sensor   string    `json:"sensor"`
    Sector   string    `json:"sector"`
    MinPressure  float32   `json:"min_pressure"`
	Coord 	 string    `json:"coord"` 
}

type Sector struct {
    Sector   string    `json:"sector"`
	Coords   string    `json:"coords"`
}

type Suscribe struct {
    Sector   Sector    `json:"Sector"`
    Queue    string    `json:"Queue"`
}


var SECTOR_NAME = os.Getenv("SECTOR_NAME")
var COORDS = os.Getenv("COORDS")

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

	var conn *amqp.Connection
	var err error

	for conn == nil {
		conn, err = amqp.Dial("amqp://" + QUEUE_HOST + ":5672/")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %v", err)
			time.Sleep(5 * time.Second) // wait 5 seconds before retrying
		}
	}
	defer conn.Close()
	log.Printf("Connected to RabbitMQ")

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

	r.PUT("/Sensor/Suscribe", func(c *gin.Context) {
		var sensor Sensor

		if err := c.ShouldBindJSON(&sensor); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, err := json.Marshal(sensor)
		if err != nil {
			log.Fatalf("Error serializing sensor to JSON: %v", err)
		}

		// Send sensor to central-server

		req, http_err := http.NewRequest("POST", "http://central-server:8080/AddSensor", bytes.NewBuffer(body))
		if http_err != nil {
			log.Fatal("Error al crear la solicitud:", http_err)
			log.Printf("ERROR HTTP")
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error al enviar la solicitud:", err)
		}
		resp.Body.Close()

		return
	})


	r.PUT("/Measurement", func(c *gin.Context) {
		var measurement Measurement
		// Deserializar el body JSON en la struct Measurement
		if err := c.ShouldBindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		body, err := json.Marshal(measurement)
		if err != nil {
			log.Fatalf("Error serializing measurement to JSON: %v", err)
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




func alertCentral(){
	for{

		sector := Sector{
			Sector: SECTOR_NAME,
			Coords: COORDS,
		}

		suscribe := Suscribe{
			Queue:  QUEUE_HOST,
			Sector: sector,
		}

		// Convertir la estructura a JSON
		jsonData, marshalerr := json.Marshal(suscribe)
		if marshalerr != nil {
			log.Fatal("Error al convertir la estructura a JSON:", marshalerr)
		}

		log.Printf("mandando http request suscribe")
		// suscribirse para obtener el nombre de la queue
		req, http_err := http.NewRequest("POST", "http://central-server:8080/Sector/Suscribe", bytes.NewBuffer(jsonData))
		if http_err != nil {
			log.Fatal("Error al crear la solicitud:", http_err)
			log.Printf("ERROR HTTP")
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error al enviar la solicitud:", err)
		}
		resp.Body.Close()

		// AGREGAR SI ES DESEADO EL MANDAR TODOS LOS DATOS GUARDADOS DE HACE 5 MINUTOS HASTA 
		// AHORA DE MONGO DB

		time.Sleep(time.Minute * 5)
	}
}
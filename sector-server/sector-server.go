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
)

type Medicion struct {
    Datetime time.Time `json:"Datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"Sensor"`
    Sector   string    `json:"Sector"`
    Presion  int       `json:"Presion"`
}

var SECTOR_NAME = os.Getenv("SECTOR_NAME")
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
	
	go getQueue()

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

func getQueue(){
	for{
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
		queueName := cleanBody
		var conn *amqp.Connection
		var err error

		for conn == nil {
			conn, err = amqp.Dial("amqp://" + queueName + "/")
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

		// AGREGAR SI ES DESEADO EL MANDAR TODOS LOS DATOS GUARDADOS DE HACE 5 MINUTOS HASTA 
		// AHORA DE MONGO DB

		time.Sleep(time.Minute * 5)
	}
}
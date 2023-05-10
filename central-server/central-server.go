package main

import (
	"encoding/json"
	"database/sql"
	"net/http"
	"time"
	"fmt"
	"log"
	"os"
	"github.com/streadway/amqp"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"bytes"
)

var DB_CONN_STRING = os.Getenv("DB_CONN_STRING")
var MIRROR_DB_SERVER_HOST = os.Getenv("MIRROR_DB_SERVER_HOST")

type Medicion struct {
    Datetime time.Time `json:"Datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"Sensor"`
    Sector   string    `json:"Sector"`
    Presion  int       `json:"Presion"`
}

var consumiendo = false

func main() {
	time.Sleep(time.Second)
	//Create DB connection
	db, err := sql.Open("postgres", DB_CONN_STRING)
	if err != nil {
		log.Fatal(err)
	}
	
	r := gin.Default()

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
		
	})

	r.GET("/Sector/Suscribe", func(c *gin.Context) {
		c.JSON(http.StatusOK, "queue:5672")
		if !consumiendo {
			consumiendo = true
			go consumer("queue:5672")
		}
	})

	r.GET("/Mediciones", func(c *gin.Context) {
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM mediciones")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las mediciones
		var mediciones []Medicion

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var medicion Medicion
			err = rows.Scan(&medicion.Datetime, &medicion.Sensor, &medicion.Sector, &medicion.Presion)
			if err != nil {
				fmt.Println("Error al preparar la sentencia SQL:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Agregar la medición al slice
			mediciones = append(mediciones, medicion)
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	
		// Devolver las mediciones como JSON
		c.JSON(http.StatusOK, mediciones)
	})
	

	r.POST("/Mediciones", func(c *gin.Context) {
		var medicion Medicion
	
		// Deserializar el body JSON en la struct Medicion
		if err := c.ShouldBindJSON(&medicion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		// Insertar en la base de datos
		fmt.Println(fmt.Sprint(medicion))
	
		// Preparar la sentencia SQL de inserción
		stmt, err := db.Prepare("INSERT INTO mediciones(datetime, sensor, sector, presion) VALUES ($1, $2, $3, $4)")
		if err != nil {
			log.Printf("Error al preparar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement de la insercion"})
			return
		}
		defer stmt.Close()
	
		// Ejecutar la sentencia SQL con los valores de la medición
		_, err = stmt.Exec(medicion.Datetime, medicion.Sensor, medicion.Sector, medicion.Presion)
		if err != nil {
			log.Printf("Error al ejecutar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo insertar la medición en la base de datos"})
			return
		}
		
		// Mandar a guardar en la base de datos réplica en un hilo.
		go sendMirror(medicion)
		
		c.Status(http.StatusOK)
		return
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func sendMirror(med Medicion) {
	clientHealth := &http.Client{}
    req, _ := http.NewRequest("GET", "http://central-server:8080/healthcheck", nil)

    for i:= 0 ; i < 3 ; i++ {
        resp, err := clientHealth.Do(req)
        if err != nil {
            fmt.Println(err)
        } else if resp.StatusCode == 200 {
			body, err := json.Marshal(med)
			if err != nil {
				log.Fatalf("Error serializing medicion to JSON: %v", err)
			}
            // Crear HTTP Request al Mirror DB Server
			req, error := http.NewRequest("POST", "http://" + MIRROR_DB_SERVER_HOST+":8080/Mediciones", bytes.NewBuffer(body) )
			if error != nil {
				log.Fatalf("Error al enviar request al Mirror DB Server: %v", error)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			// Enviar la solicitud y capturar la respuesta y el posible error
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error al pegar al endpoint de mediciones: %v", err)
			}
			defer resp.Body.Close()
            break
        }
        time.Sleep(time.Minute * 1) // espera 1 Minuto antes de volver a intentar
    }
	
}
// Thread
func consumer(queue string) {

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
			// Crear una solicitud HTTP POST con el cuerpo JSON y el encabezado "Content-Type: application/json"
			body := bytes.NewBuffer(d.Body)

			req, err := http.NewRequest("POST", "http://central-server:8080/Mediciones", body)
			req.Header.Set("Content-Type", "application/json")

			// Enviar la solicitud y capturar la respuesta y el posible error
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error al pegar al endpoint de mediciones: %v", err)
			}
			defer resp.Body.Close()

			// Imprimir el código de estado de la respuesta
			fmt.Println(resp.StatusCode)
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
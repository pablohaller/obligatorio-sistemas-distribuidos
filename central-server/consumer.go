// consumer.go

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

func main() {
	time.Sleep(time.Second)
	//Create DB connection
	db, err := sql.Open("postgres", DB_CONN_STRING)
	if err != nil {
		log.Fatal(err)
	}
	
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong!")
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
	

	r.POST("/Mediciones", func(c *gin.Context){
		var medicion Medicion

        // Deserializar el body JSON en la struct Medicion
        if err := c.ShouldBindJSON(&medicion); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        // Insertar en la base de datos
        fmt.Println(fmt.Sprint(medicion))
		// Preparar la sentencia SQL de inserción
		stmt, err := db.Prepare("INSERT INTO mediciones(datetime, sensor, sector, presion) VALUES($1, $2, $3, $4)")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// Ejecutar la sentencia SQL con los valores de la medición
		_, err = stmt.Exec(medicion.Datetime, medicion.Sensor, medicion.Sector, medicion.Presion)
		if err != nil {
			fmt.Println("Error al ejecutar la sentencia SQL:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
			
        c.Status(http.StatusOK)
		return 
	})


	r.Run() // listen and serve on 0.0.0.0:8080
}

// Thread
func consumer() {
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
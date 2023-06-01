package main

import (
	//"encoding/json"
	"bytes"
	"database/sql"
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

var DB_CONN_STRING = os.Getenv("DB_CONN_STRING")
var MIRROR_DB_SERVER_HOST = os.Getenv("MIRROR_DB_SERVER_HOST")

type Medicion struct {
    Datetime time.Time `json:"Datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"Sensor"`
    Sector   string    `json:"Sector"`
    Presion  int       `json:"Presion"`
}

type Suscribe struct {
    Sector   string    `json:"Sector"`
    Queue    string    `json:"Queue"`
}

var consumiendo []string

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

	r.POST("/Sector/Suscribe", func(c *gin.Context) {
		var found bool
		var suscribe Suscribe
		if err := c.ShouldBindJSON(&suscribe); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Recorrer el array de strings
		for _, str := range consumiendo {
			// Comparar el elemento actual con el valor deseado
			if str == suscribe.Sector {
				// Se encontró el valor deseado
				found = true
				break
			}
		}
		if !found {
			go consumer(suscribe.Queue+ ":5672")
			consumiendo = append(consumiendo, suscribe.Sector)
			c.JSON(http.StatusCreated, suscribe.Queue)
			log.Printf(suscribe.Sector + " se suscribio para ser consumido")
		}else{
			c.JSON(http.StatusOK, "Already consuming")
			log.Printf(suscribe.Sector + " ya esta suscripto")
		}
	})

	// Definir el endpoint "/Mediciones" con un parámetro de ruta para los minutos
	r.GET("/UltMediciones/:minutos", func(c *gin.Context) {
		// Obtener el valor de los minutos desde la URL
		minutosStr := c.Param("minutos")
		minutos, err := strconv.Atoi(minutosStr)
		if err != nil {
			// Manejar el error si no se puede convertir a entero
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido para minutos"})
			return
		}

		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM mediciones WHERE \"datetime\" >= NOW() - INTERVAL '" + strconv.Itoa(minutos) + " minutes'")
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
	

	r.PUT("/Mediciones", func(c *gin.Context) {
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
	
		c.Status(http.StatusOK)
		return
	})

	r.GET("/Alerta", func(c *gin.Context) {
		var medicion Medicion
	
		// Deserializar el body JSON en la struct Medicion
		if err := c.ShouldBindJSON(&medicion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM mediciones WHERE \"datetime\" >= (SELECT NOW() - INTERVAL '3 hours') - INTERVAL '" + strconv.Itoa(minutos) + " minutes'")
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
		return
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
	
// Thread
func consumer(queue string) {

	var conn *amqp.Connection
	var err error
	log.Printf(queue + " -- Dentro del consumer")
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
			// Crear una solicitud HTTP PUT con el cuerpo JSON y el encabezado "Content-Type: application/json"
			body := bytes.NewBuffer(d.Body)

			req, err := http.NewRequest("PUT", "http://central-server:8080/Mediciones", body)
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
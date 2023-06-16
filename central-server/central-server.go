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

type Measurement struct {
    Datetime time.Time `json:"datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"sensor"`
    Sector   string    `json:"sector"`
    Pressure  float32       `json:"pressure"`
}

type Suscribe struct {
    Sector   Sector    `json:"sector"`
    Queue    string    `json:"queue"`
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

	r.POST("/AddSensor", func(c *gin.Context) {
		var sensor Sensor
		if err := c.ShouldBindJSON(&sensor); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Add sensor to central-db

		fmt.Println(fmt.Sprint(sensor))

		var found bool
		// Recorrer el array de strings
		for _, str := range consumiendo {
			// Comparar el elemento actual con el valor deseado
			if str == sensor.Sector {
				// Se encontró el valor deseado
				found = true
				break
			}
		}
		
		if found {
			// Preparar la sentencia SQL de inserción
			stmt, err := db.Prepare("INSERT INTO sensors (sensor, sector, min_pressure, coord) VALUES ($1, $2, $3, $4)")
			if err != nil {
				log.Printf("Error al preparar la sentencia SQL: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement de la insercion"})
				return
			}
			defer stmt.Close()
		
			// Ejecutar la sentencia SQL con los valores de la medición
			_, err = stmt.Exec(sensor.Sensor, sensor.Sector, sensor.MinPressure, sensor.Coord)
			if err != nil {
				log.Printf("Error al ejecutar la sentencia SQL: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo insertar la medición en la base de datos"})
				return
			}
		
			c.Status(http.StatusOK)
		}else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "El sector no se encuentra registrado"})
		}
	})

	r.POST("/Sector/Suscribe", func(c *gin.Context) {
		var found bool
		var suscribe Suscribe
		log.Printf("Starting /Sector/Suscribe:", suscribe.Sector.Sector)
		if err := c.ShouldBindJSON(&suscribe); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Recorrer el array de strings
		for _, str := range consumiendo {
			// Comparar el elemento actual con el valor deseado
			if str == suscribe.Sector.Sector {
				// Se encontró el valor deseado
				found = true
				break
			}
		}
		if !found {
			go consumer(suscribe.Queue+ ":5672")
			consumiendo = append(consumiendo, suscribe.Sector.Sector)

			// Add sector to central-db
		
			fmt.Println(fmt.Sprint(suscribe.Sector))
		
			// Preparar la sentencia SQL de inserción
			stmt, err := db.Prepare("INSERT INTO sectors (sector, coords) VALUES ($1, $2)")
			if err != nil {
				log.Printf("Error al preparar la sentencia SQL: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement de la insercion"})
				return
			}
			defer stmt.Close()
		
			// Ejecutar la sentencia SQL con los valores de la medición
			_, err = stmt.Exec(suscribe.Sector.Sector,suscribe.Sector.Coords)
			if err != nil {
				log.Printf("Error al ejecutar la sentencia SQL: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo insertar la medición en la base de datos"})
				return
			}
			log.Println("Se ha ejecutado la insercion de sector de forma exitosa.")
		
			c.Status(http.StatusOK)
			return

			//////////////////////////////

			c.JSON(http.StatusCreated, suscribe.Queue)
			log.Printf(suscribe.Sector.Sector + " se suscribio para ser consumido")
		}else{
			c.JSON(http.StatusOK, "Already consuming")
			log.Printf(suscribe.Sector.Sector + " ya esta suscripto")
		}
	})

	// Definir el endpoint "/Measurement" con un parámetro de ruta para los minutes
	r.GET("/LastMeasurements/:minutes", func(c *gin.Context) {
		// Obtener el valor de los minutes desde la URL
		minutesStr := c.Param("minutes")
		minutes, err := strconv.Atoi(minutesStr)
		if err != nil {
			// Manejar el error si no se puede convertir a entero
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido para minutes"})
			return
		}

		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM measurements WHERE \"datetime\" >= (SELECT NOW() - INTERVAL '3 hours') - INTERVAL '" + strconv.Itoa(minutes) + " minutes'")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las measurements
		var measurements []Measurement

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var measurement Measurement
			err = rows.Scan(&measurement.Datetime, &measurement.Sensor, &measurement.Sector, &measurement.Pressure)
			if err != nil {
				fmt.Println("Error al preparar la sentencia SQL:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Agregar la medición al slice
			measurements = append(measurements, measurement)
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	
		// Devolver las mediciones como JSON
		c.JSON(http.StatusOK, measurements)
	})

	r.GET("/Measurement", func(c *gin.Context) {
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM measurements")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las mediciones
		var measurements []Measurement

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var measurement Measurement
			err = rows.Scan(&measurement.Datetime, &measurement.Sensor, &measurement.Sector, &measurement.Pressure)
			if err != nil {
				fmt.Println("Error al preparar la sentencia SQL:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Agregar la medición al slice
			measurements = append(measurements, measurement)
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	
		// Devolver las mediciones como JSON
		c.JSON(http.StatusOK, measurements)
	})
	

	r.PUT("/Measurements", func(c *gin.Context) {
		var measurement Measurement
	
		// Deserializar el body JSON en la struct Medicion
		if err := c.ShouldBindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		// Insertar en la base de datos
		fmt.Println(fmt.Sprint(measurement))
	
		// Preparar la sentencia SQL de inserción
		stmt, err := db.Prepare("INSERT INTO measurements (datetime, sensor, sector, pressure) VALUES ($1, $2, $3, $4)")
		if err != nil {
			log.Printf("Error al preparar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement de la insercion"})
			return
		}
		defer stmt.Close()
	
		// Ejecutar la sentencia SQL con los valores de la medición
		_, err = stmt.Exec(measurement.Datetime, measurement.Sensor, measurement.Sector, measurement.Pressure)
		if err != nil {
			log.Printf("Error al ejecutar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo insertar la medición en la base de datos"})
			return
		}
	
		c.Status(http.StatusOK)
		return
	})

	r.GET("/Alert", func(c *gin.Context) {
		var measurement Measurement
	
		// Deserializar el body JSON en la struct Medicion
		if err := c.ShouldBindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM Mediciones WHERE datetime >= '"+ measurement.Datetime.Format("2006-01-02T15:04:05Z") +"' AND sensor = '"+measurement.Sensor+"' AND sector = '"+measurement.Sector+"';");
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las mediciones
		var measurements []Measurement

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var measurement Measurement
			err = rows.Scan(&measurement.Datetime, &measurement.Sensor, &measurement.Sector, &measurement.Pressure)
			if err != nil {
				fmt.Println("Error al preparar la sentencia SQL:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Agregar la medición al slice
			measurements = append(measurements, measurement)
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	
		// Devolver las mediciones como JSON
		c.JSON(http.StatusOK, measurements)
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

			req, err := http.NewRequest("PUT", "http://central-server:8080/Measurements", body)
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
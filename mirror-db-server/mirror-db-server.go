package main

import (
	_ "encoding/json"
	"database/sql"
	"net/http"
	"time"
	"fmt"
	"log"
	"os"
	"strconv"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var DB_CONN_STRING = os.Getenv("DB_CONN_STRING")

type Measurement struct {
    Datetime time.Time `json:"datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"sensor"`
    Sector   string    `json:"sector"`
    Pressure  int       `json:"pressure"`
}

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

	// Definir el endpoint "/Measurements" con un parámetro de ruta para los minutes
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

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Measurement
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
	
		// Devolver las measurements como JSON
		c.JSON(http.StatusOK, measurements)
	})

	r.GET("/Alert", func(c *gin.Context) {
		var measurement Measurement
	
		// Deserializar el body JSON en la struct Measurement
		if err := c.ShouldBindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM Measurementes WHERE datetime >= '"+ measurement.Datetime.Format("2006-01-02T15:04:05Z") +"' AND sensor = '"+measurement.Sensor+"' AND sector = '"+measurement.Sector+"';");
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las measurements
		var measurements []Measurement

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Measurement
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
	
		// Devolver las measurements como JSON
		c.JSON(http.StatusOK, measurements)
		return
	})


	r.GET("/Measurementes", func(c *gin.Context) {
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM measurements")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las measurements
		var measurements []Measurement

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Measurement
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
	
		// Devolver las measurements como JSON
		c.JSON(http.StatusOK, measurements)
	})

	
	
	r.Run() // listen and serve on 0.0.0.0:8080
}
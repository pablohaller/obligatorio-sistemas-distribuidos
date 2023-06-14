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

type Medicion struct {
    Datetime time.Time `json:"Datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"Sensor"`
    Sector   string    `json:"Sector"`
    Presion  int       `json:"Presion"`
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
	})

	r.GET("/Alerta", func(c *gin.Context) {
		var medicion Medicion
	
		// Deserializar el body JSON en la struct Medicion
		if err := c.ShouldBindJSON(&medicion); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM Mediciones WHERE datetime >= '"+ medicion.Datetime.Format("2006-01-02T15:04:05Z") +"' AND sensor = '"+medicion.Sensor+"' AND sector = '"+medicion.Sector+"';");
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

	
	
	r.Run() // listen and serve on 0.0.0.0:8080
}
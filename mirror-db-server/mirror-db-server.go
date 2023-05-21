package main

import (
	_ "encoding/json"
	"database/sql"
	"net/http"
	"time"
	"fmt"
	"log"
	"os"
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
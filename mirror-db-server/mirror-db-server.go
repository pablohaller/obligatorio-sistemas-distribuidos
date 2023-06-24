package main

import (
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
)

var DB_CONN_STRING = os.Getenv("DB_CONN_STRING")


type Measurement struct {
    Datetime time.Time `json:"datetime" time_format:"2006-01-02 15:04:05"`
    Sensor   string    `json:"sensor"`
    Sector   string    `json:"sector"`
    Pressure  float32       `json:"pressure"`
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

type MeasurementData struct {
	Datetime time.Time 		`json:"datetime" time_format:"2006-01-02 15:04:05"`
	Pressure  float32       `json:"pressure"`
}

type SectorMap struct {
	Coords   string    		`json:"coords"`
	Sensors []SensorMap	`json:"sensors"`
}

type SensorMap struct {
	Sensor   string    `json:"sensor"`
	Coord 	 string    `json:"coord"` 
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

	r.Use(corsMiddleware())

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
		
	})

	// Agregar en mirror
	r.GET("/Sectors", func(c *gin.Context) {
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM sectors")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las mediciones
		var sectors []Sector

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var sector Sector
			err = rows.Scan(&sector.Sector, &sector.Coords)
			if err != nil {
				fmt.Println("Error al leer la consulta:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Agregar la medición al slice
			sectors = append(sectors, sector)
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	
		// Devolver las sectores como JSON
		c.JSON(http.StatusOK, sectors)
	})

	// Agregar en mirror
	r.GET("/Sensors", func(c *gin.Context) {
		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM sensors")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las mediciones
		var sensors []Sensor

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var sensor Sensor
			err = rows.Scan(&sensor.Sensor, &sensor.Sector, &sensor.MinPressure, &sensor.Coord)
			if err != nil {
				fmt.Println("Error al leer la consulta:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Agregar la medición al slice
			sensors = append(sensors, sensor)
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	
		// Devolver las sensores como JSON
		c.JSON(http.StatusOK, sensors)
	})

	r.POST("/MapReport", func(c *gin.Context){

		var measurement Measurement
		
		if err := c.ShouldBindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.Prepare("SELECT * FROM sectors WHERE \"sector\" = $1")
		if err != nil {
			log.Printf("Error al preparar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement del select"})
			return
		}
		defer stmt.Close()

		// Ejecutar la sentencia SQL con los valores de la medición
		rows, err := stmt.Query(measurement.Sector)
		if err != nil {
			log.Printf("Error al ejecutar la sentencia SQL: %v", err)
			return
		}
		defer rows.Close() // remember to close the rows object when done 

		var sectorMap SectorMap
		for rows.Next() {
			var sector Sector
			err = rows.Scan(&sector.Sector, &sector.Coords)
			if err != nil {
				fmt.Println("Error al leer la consulta:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			sectorMap = SectorMap{
				Coords: sector.Coords,
				Sensors: []SensorMap{},
			}
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		stmt, err = db.Prepare("SELECT * FROM sensors WHERE \"sensor\" = $1")
		if err != nil {
			log.Printf("Error al preparar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement del select"})
			return
		}
		defer stmt.Close()

		// Ejecutar la sentencia SQL con los valores de la medición
		rows, err = stmt.Query(measurement.Sensor)
		if err != nil {
			log.Printf("Error al ejecutar la sentencia SQL: %v", err)
			return
		}

		for rows.Next() {
			var sensor Sensor
			err = rows.Scan(&sensor.Sensor, &sensor.Sector,&sensor.MinPressure,&sensor.Coord)
			if err != nil {
				fmt.Println("Error al leer la consulta:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			
			sectorMap.Sensors = append(sectorMap.Sensors, SensorMap{
				Sensor: sensor.Sensor,
				Coord : sensor.Coord,
			})
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		
		c.JSON(http.StatusOK, sectorMap)
	})

	// Agregar en mirror
	r.GET("/Map", func(c *gin.Context) {
		//Para la lectura de la base de datos:
		
		rows, err := db.Query("SELECT * FROM sectors ")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Slice para guardar las mediciones
		dict := map[string]SectorMap{}
		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var sector Sector
			err = rows.Scan(&sector.Sector, &sector.Coords)
			if err != nil {
				fmt.Println("Error al leer la consulta:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			// Agregar la medición al slice

			dict[sector.Sector] = SectorMap{
				Coords : sector.Coords,
				Sensors: []SensorMap{},
			}
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		//Para la lectura de la base de datos:
		rows, err = db.Query("SELECT * FROM sensors")
		if err != nil {
			fmt.Println("Error al preparar la sentencia SQL:", err)
		}
		defer rows.Close() // remember to close the rows object when done

		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var sensor Sensor
			err = rows.Scan(&sensor.Sensor, &sensor.Sector, &sensor.MinPressure, &sensor.Coord)
			if err != nil {
				fmt.Println("Error al leer la consulta:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			
			SensorMap := SensorMap{
				Sensor: sensor.Sensor,
				Coord: sensor.Coord,
			}

			sensors := append(dict[sensor.Sector].Sensors, SensorMap)

			dict[sensor.Sector] = SectorMap{
				Coords: dict[sensor.Sector].Coords,
				Sensors: sensors,
			}
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		
		// Devolver el diccionario como JSON
		c.JSON(http.StatusOK, dict)
	})
	

	// Agregar en mirror
	r.GET("LastSectorMeasurements/:sector/:minutes", func(c *gin.Context) {
		// Obtener el valor de los minutes desde la URL
		minutesStr := c.Param("minutes")
		sectorStr := c.Param("sector")
		_, err := strconv.Atoi(minutesStr)
		if err != nil {
			// Manejar el error si no se puede convertir a entero
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido para minutes"})
			return
		}

		// Preparar la sentencia: SELECT * FROM measurements WHERE sector = 'sector-server-a' and datetime >= (SELECT NOW() - INTERVAL '3 hours') - INTERVAL '2 minutes'

		stmt, err := db.Prepare("SELECT * FROM measurements WHERE \"sector\" = $1 and \"datetime\" >= (SELECT NOW() - INTERVAL '3 hours') - INTERVAL '" + minutesStr + " minutes'")
		if err != nil {
			log.Printf("Error al preparar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement del select"})
			return
		}
		defer stmt.Close()

		// Ejecutar la sentencia SQL con los valores de la medición
		rows, err := stmt.Query(sectorStr)
		if err != nil {
			log.Printf("Error al ejecutar la sentencia SQL: %v", err)
			return
		}
		defer rows.Close() // remember to close the rows object when done 

		dict := map[string][]MeasurementData{}
		// Recorrer el resultado de la consulta y guardar los valores en las estructuras de tipo Medicion
		for rows.Next() {
			var measurement Measurement
			err = rows.Scan(&measurement.Datetime, &measurement.Sensor, &measurement.Sector, &measurement.Pressure)
			if err != nil {
				fmt.Println("Error al leer la consulta:", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			
			md := MeasurementData{
				Datetime: measurement.Datetime,
				Pressure: measurement.Pressure,
			}

			value, ok := dict[measurement.Sensor]
			if ok {
				dict[measurement.Sensor] = append(value, md)
			} else {
				list := []MeasurementData{}
				list = append(list, md)
				dict[measurement.Sensor] = list
			}
		}
		// Si hubo algún error al recorrer los resultados
		if err = rows.Err(); err != nil {
			fmt.Println("Error al recorrer los resultados de la consulta:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		
		c.JSON(http.StatusOK, dict)
	})

	// Definir el endpoint "/Measurement" con un parámetro de ruta para los minutes
	r.GET("/LastMeasurements/:minutes", func(c *gin.Context) {
		// Obtener el valor de los minutes desde la URL
		minutesStr := c.Param("minutes")
		_, err := strconv.Atoi(minutesStr)
		if err != nil {
			// Manejar el error si no se puede convertir a entero
			c.JSON(http.StatusBadRequest, gin.H{"error": "Valor inválido para minutes"})
			return
		}

		//Para la lectura de la base de datos:
		rows, err := db.Query("SELECT * FROM measurements WHERE \"datetime\" >= (SELECT NOW() - INTERVAL '3 hours') - INTERVAL '" + minutesStr + " minutes'")
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
				fmt.Println("Error al leer la consulta:", err)
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

	r.GET("/Measurements", func(c *gin.Context) {
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
				fmt.Println("Error al leer la consulta:", err)
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

	r.POST("/Alert", func(c *gin.Context) {
		var measurement Measurement
	
		// Deserializar el body JSON en la struct Medicion
		if err := c.ShouldBindJSON(&measurement); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		stmt, err := db.Prepare("SELECT * FROM measurements WHERE datetime >= $1 AND sensor = $2 AND sector = $3")
		if err != nil {
			log.Printf("Error al preparar la sentencia SQL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo preparar el statement del select"})
			return
		}
		defer stmt.Close()

		//Para la lectura de la base de datos:
		rows, err := stmt.Query(measurement.Datetime.Format("2006-01-02T15:04:05Z"),measurement.Sensor, measurement.Sector);
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
				fmt.Println("Error al leer la consulta:", err)
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

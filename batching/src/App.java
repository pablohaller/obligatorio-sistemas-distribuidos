package batching;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Random;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.*;
import com.sun.net.httpserver.HttpServer;
import java.net.InetSocketAddress;
import java.net.URI;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.io.BufferedReader;
import java.io.IOException;

import java.io.IOException;
import java.io.OutputStream;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;


public class App {  
    // Hashmap 
    private static HashMap<String, String> reportedSensors = new HashMap<>();
    private static int pollingRate;
    private static HashMap<String,Sensor> sensors = new HashMap<>();
    private static String webServer = System.getenv("WEB-SERVER-HOST");
    public static void main(String[] args) throws Exception {
        pollingRate =  Integer.valueOf(System.getenv("POLLING-RATE"));
        while (true) {
            List<Measurement> alerts = getAlerts();
            if(!alerts.isEmpty()){
                int initBackoff = 2;
                int maxBackoff = 20;
                int retries = 0;
                Random r = new Random();
                while (!alerts.isEmpty()) {
                    URL url = new URL("http://" + webServer + ":3000/api/measures");
                    System.out.println(url);

                    // Open a connection to the URL
                    HttpURLConnection connection = (HttpURLConnection) url.openConnection();

                    // Set the request method (GET, POST, PUT, DELETE, etc.)
                    connection.setRequestMethod("POST");

                    // Optional: Set request headers
                    connection.setRequestProperty("Content-Type", "application/json");

                    connection.setDoOutput(true);
                    OutputStream outputStream = connection.getOutputStream();
                    String data = "{\"data\":\"" + MeasurementToJson(alerts.get(0)).replace("\"","'") + "\", \"filtration\":true}"; // Intentar agregar elemento
                    // 

                    System.out.println("JSON request:\n" + data);
                    outputStream.write(data.getBytes("UTF-8"));
                    outputStream.flush();
                    outputStream.close();          

                    // Get the response code
                    int responseCode = connection.getResponseCode();

                    if (responseCode == 200) {
                        // Read the response body
                        BufferedReader reader = new BufferedReader(new InputStreamReader(connection.getInputStream()));
                        String line;
                        StringBuilder response = new StringBuilder();
                        while ((line = reader.readLine()) != null) {
                            response.append(line);
                        }
                        reader.close();

                        String report = response.toString();

                        // Crear una instancia del ObjectMapper
                        ObjectMapper objectMapper = new ObjectMapper();

                        // Analizar el JSON y obtener el nodo raíz
                        JsonNode rootNode = objectMapper.readTree(report);

                        // Obtener el valor del campo "id" como cadena (String)
                        String id = rootNode.get("id").asText();

                        // Imprimir el valor del campo "id"
                        System.out.println("ID: " + id);

                        // Asociar el codigo del reporte con el sensor para despues en el filter poder darlo de baja si es necesario
                        reportedSensors.put(alerts.get(0).sector + alerts.get(0).sensor, id);

                        // Sacar el elemento que acabas de guardar.
                        alerts.remove(0); 
                        retries = 0;
                    }
                    
                    System.out.println("Response Code: " + responseCode);
                    retries++;
                    double random = Math.min(maxBackoff, initBackoff * Math.pow(2, retries));
                    Thread.sleep(1000 * r.nextInt((int)random + 1));
                }
            }

            Thread.sleep(1000*pollingRate*60);
        }
    }

    public static List<Measurement> getAlerts() {
        String json = "";
        try {
            // Get the value of the environment variable
            String mirrorDbServer = System.getenv("MIRROR-DB-SERVER");

            // Create a URL object with the endpoint you want to send the request to
            URL url = new URL("http://" + mirrorDbServer + ":8080/LastMeasurements/" + pollingRate);
            System.out.println(url);

            // Open a connection to the URL
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();

            // Set the request method (GET, POST, PUT, DELETE, etc.)
            connection.setRequestMethod("GET");

            // Optional: Set request headers
            connection.setRequestProperty("Content-Type", "application/json");

            // Get the response code
            int responseCode = connection.getResponseCode();
            System.out.println("Response Code: " + responseCode);

            // Read the response
            BufferedReader reader = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            String line;
            StringBuilder response = new StringBuilder();
            while ((line = reader.readLine()) != null) {
                response.append(line);
            }
            reader.close();

            // Print the response
            System.out.println("Response: " + response.toString());
            json = response.toString();

            // {{},{},{},{}}
            // Close the connection
            connection.disconnect();

            
        } catch (IOException e) {
            e.printStackTrace();
        }
        if(json != ""){
            return filter(JsonToMeasurements(json));
        }
        return null;
    }

    public static HashMap<String, Sensor> getSensors() {
        String json = "";
        try {
            // Get the value of the environment variable
            String mirrorDbServer = System.getenv("MIRROR-DB-SERVER");

            // Create a URL object with the endpoint you want to send the request to
            URL url = new URL("http://" + mirrorDbServer + ":8080/Sensors");
            System.out.println(url);

            // Open a connection to the URL
            HttpURLConnection connection = (HttpURLConnection) url.openConnection();

            // Set the request method (GET, POST, PUT, DELETE, etc.)
            connection.setRequestMethod("GET");

            // Optional: Set request headers
            connection.setRequestProperty("Content-Type", "application/json");

            // Get the response code
            int responseCode = connection.getResponseCode();
            System.out.println("Response Code: " + responseCode);

            // Read the response
            BufferedReader reader = new BufferedReader(new InputStreamReader(connection.getInputStream()));
            String line;
            StringBuilder response = new StringBuilder();
            while ((line = reader.readLine()) != null) {
                response.append(line);
            }
            reader.close();

            // Print the response
            System.out.println("Sensors: " + response.toString());
            json = response.toString();

            // {{},{},{},{}}
            // Close the connection
            connection.disconnect();

            
        } catch (IOException e) {
            e.printStackTrace();
        }
        if(json != ""){
            return JsonToSensors(json);
        }
        return new HashMap<String, Sensor>();
    }

    public static HashMap<String,Sensor> JsonToSensors(String jsonSensors){
        HashMap<String, Sensor> sensorsMap = new HashMap<>();
        ObjectMapper objectMapper = new ObjectMapper();
        // Parse JSON string into a JsonNode array
        JsonNode jsonNodeArray;
        Sensor sensor;
        try {
            jsonNodeArray = objectMapper.readTree(jsonSensors);
            for (JsonNode jsonNode : jsonNodeArray) {
                // Convert individual object to a JSON string
                String jsonSensor = jsonNode.toString();
                // Print the individual JSON string
                sensor = objectMapper.readValue(jsonSensor, Sensor.class);
                sensorsMap.put(sensor.sector + sensor.sensor, sensor);
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        // Iterate over each object in the array    
        return sensorsMap;
    }

    public static List<Measurement> JsonToMeasurements(String jsonMeasurements){
        ObjectMapper objectMapper = new ObjectMapper();
        // Parse JSON string into a JsonNode array
        JsonNode jsonNodeArray;
        List<Measurement> measurements = new LinkedList<>();
        try {
            jsonNodeArray = objectMapper.readTree(jsonMeasurements);
            for (JsonNode jsonNode : jsonNodeArray) {
                // Convert individual object to a JSON string
                String jsonMeasurement = jsonNode.toString();
                // Print the individual JSON string
                measurements.add(objectMapper.readValue(jsonMeasurement, Measurement.class));
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        // Iterate over each object in the array    
        return measurements;
    }

    public static String MeasurementsToJson(List<Measurement> measurements){
        if (measurements.size() > 0){
            // Crear una instancia de ObjectMapper
            ObjectMapper objectMapper = new ObjectMapper();
            try {
                // Convertir la lista de objetos a JSON
                String json = objectMapper.writeValueAsString(measurements);
                return json;
            } catch (JsonProcessingException e) {
                e.printStackTrace();
            }
        }
        return "";
    }

    public static String MeasurementToJson(Measurement measurement){
        // Crear una instancia de ObjectMapper
        ObjectMapper objectMapper = new ObjectMapper();
        try {
            // Convertir el objeto a JSON
            String json = objectMapper.writeValueAsString(measurement);
            return json;
        } catch (JsonProcessingException e) {
            e.printStackTrace();
        }
        return "";
    }   


    public static List<Measurement> filter(List <Measurement> measurements) {
        List<Measurement> filterList = new LinkedList<>();

        for (Measurement m : measurements) {
            
            Sensor sen = sensors.get(m.sector + m.sensor);
            if (sen == null) {
                sensors = getSensors();
                System.out.println("Hashmap size after getSensors(): " + sensors.size());
                sen = sensors.get(m.sector + m.sensor);
            }
            if (m.getPressure() < sen.getMinPressure() && !reportedSensors.containsKey(m.sector + m.sensor)) { // Si la presion es menor al umbral y el sensor no había sIdo procesado.
                filterList.add(m);
                reportedSensors.put(m.sector + m.sensor,"");
                System.out.println("Report added: "+m.sector + m.sensor);
            }else if(m.getPressure() > sen.getMinPressure() && reportedSensors.containsKey(m.sector + m.sensor)  ) {
                //PEGARLE AL ENDPOINT QUE HACE EL WEBSERVER
                System.out.println("else if joined, sending id: " + reportedSensors.get(m.sector + m.sensor));
                sendRemoveReport(reportedSensors.get(m.sector + m.sensor));
                reportedSensors.remove(m.sector + m.sensor);
                System.out.println("La key sigue estando? "+ reportedSensors.containsKey(m.sector + m.sensor));
            
            }
        }
        System.out.println("Measurements after filtering: "+filterList.size());
        return filterList;
    }


    public static void sendRemoveReport(String id){ 
        int initBackoff = 2;
        int maxBackoff = 20;
        int retries = 0;
        Random r = new Random();
        System.out.println("New Call id: " + id);
        while (true) {
            try {
                //URL url = new URL("http://" + webServer + ":3000/api/measures/?id=" + id);
                URL url = new URL("http://" + webServer + ":3000/api/measures?id=" + id);

                System.out.println(url);

                // Open a connection to the URL
                HttpURLConnection connection = (HttpURLConnection) url.openConnection();

                // Set the request method to POST
                connection.setRequestMethod("DELETE");

                // Optional: Set request headers
                connection.setRequestProperty("Content-Type", "application/json");
                
                // Get the response code
                int responseCode = connection.getResponseCode(); 
                System.out.println("Response Code: " + responseCode);
                if (responseCode == 200) {
                    break;
                }
                
                retries++;
                double random = Math.min(maxBackoff, initBackoff * Math.pow(2, retries));
                Thread.sleep(1000 * r.nextInt((int)random + 1));
            } catch (Exception ex) {
                ex.printStackTrace();
            }
            
        }
    }

    public static class Sensor {

        private String sensor;         
        private String sector;     
        private float minPressure;
        private String coord;
        
        public Sensor (){};

        public Sensor (String sensor, String sector, float minPressure, String coord) {
            this.sensor = sensor;
            this.sector = sector;
            this.minPressure = minPressure;
            this.coord = coord;
        }
    
        @JsonProperty("sensor")
        public void setSensor(String sensor) {
            this.sensor = sensor;
        }

        @JsonProperty("sensor")
        public String getSensor() {
            return sensor;
        }
    
        @JsonProperty("sector")
        public void setSector(String sector) {
            this.sector = sector;
        }

        @JsonProperty("sector")
        public String getSector() {
            return sector;
        }
    
        @JsonProperty("min_pressure")
        public void setMinPressure(float pressure) {
            this.minPressure = pressure;
        }

        @JsonProperty("min_pressure")
        public float getMinPressure() {
            return minPressure;
        }

        @JsonProperty("coord")
        public void setCoord(String coord) {
            this.coord = coord;
        }

        @JsonProperty("coord")
        public String getCoord() {
            return coord;
        }
    }

    public static class Measurement {

        private String datetime;
        private String sensor;  
        private String sector;
        private float pressure;

        public Measurement (){};

        public Measurement (String datetime, String sensor, String sector, float pressure) {
            this.datetime = datetime;
            this.sensor = sensor;
            this.sector = sector;
            this.pressure = pressure;
        }
    
        @JsonProperty("datetime")
        public void setDatetime(String datetime) {
            this.datetime = datetime;
        }
        
        @JsonProperty("datetime")
        public String getDatetime() {
            return datetime;
        }
    
        @JsonProperty("sensor")
        public void setSensor(String sensor) {
            this.sensor = sensor;
        }

        @JsonProperty("sensor")
        public String getSensor() {
            return sensor;
        }
    
        @JsonProperty("sector")
        public void setSector(String sector) {
            this.sector = sector;
        }

        @JsonProperty("sector")
        public String getSector() {
            return sector;
        }
    
        @JsonProperty("pressure")
        public void setPressure(float pressure) {
            this.pressure = pressure;
        }

        @JsonProperty("pressure")
        public float getPressure() {
            return pressure;
        }
    }
}


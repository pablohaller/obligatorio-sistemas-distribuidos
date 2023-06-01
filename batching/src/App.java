package batching;

import java.util.ArrayList;
import java.util.LinkedList;
import java.util.List;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.*;
import com.sun.net.httpserver.HttpServer;
import java.net.InetSocketAddress;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.io.BufferedReader;
import java.io.IOException;

import java.io.IOException;
import java.io.OutputStream;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;


public class App {  
    private static List<String> sensoresEnFuga = new ArrayList<>();
    private static int pollingRate;
    public static void main(String[] args) throws Exception {
        pollingRate =  Integer.valueOf(System.getenv("POLLING-RATE"));
        String webServer = System.getenv("WEB-SERVER-HOST");
        /* int port = 8080;
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        server.createContext("/Batching", new MyHttpHandler());
        server.setExecutor(null); // Utiliza el executor por defecto
        server.start();
        System.out.println("Servidor iniciado en el puerto " + port);  */
        while (true) {

            String alertas = getAlertas();
            if(alertas != ""){
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
                String data = "{\"data\":\"" + alertas.replace("\"","'") + "\", \"filtration\":true}";

                System.out.println("JSON QUE SE CARGA EN LA REQUEST:\n" + data);
                outputStream.write(data.getBytes("UTF-8"));
                outputStream.flush();
                outputStream.close();          

                // Get the response code
                int responseCode = connection.getResponseCode();
                System.out.println("Response Code: " + responseCode);
            }

            Thread.sleep(1000*pollingRate*60);
        }
    }

    public static String getAlertas() {
        String json = "";
        try {
            // Get the value of the environment variable
            String mirrorDbServer = System.getenv("MIRROR-DB-SERVER");

            // Create a URL object with the endpoint you want to send the request to
            URL url = new URL("http://" + mirrorDbServer + ":8080/UltMediciones/" + pollingRate);
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
            return MedicionesToJson(filter(JsonToMediciones(json)));
        }
        return "";
    }

    public static List<Medicion> JsonToMediciones(String jsonMediciones){
        ObjectMapper objectMapper = new ObjectMapper();
        // Parse JSON string into a JsonNode array
        JsonNode jsonNodeArray;
        List<Medicion> mediciones = new LinkedList<>();
        try {
            jsonNodeArray = objectMapper.readTree(jsonMediciones);
            for (JsonNode jsonNode : jsonNodeArray) {
                // Convert individual object to a JSON string
                String jsonMedicion = jsonNode.toString();
                // Print the individual JSON string
                mediciones.add(objectMapper.readValue(jsonMedicion, Medicion.class));
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
        // Iterate over each object in the array    
        return mediciones;
    }

    public static String MedicionesToJson(List<Medicion> mediciones){
        if (mediciones.size() > 0){
            // Crear una instancia de ObjectMapper
            ObjectMapper objectMapper = new ObjectMapper();
            try {
                // Convertir la lista de objetos a JSON
                String json = objectMapper.writeValueAsString(mediciones);
                return json;
            } catch (JsonProcessingException e) {
                e.printStackTrace();
            }
        }
        return "";
    }

    public static List<Medicion> filter(List <Medicion> mediciones) {
        List<Medicion> filterList = new LinkedList<>();
        System.out.println("cantidad de mediciones: " + mediciones.size());
        for (Medicion m : mediciones) {
            if (m.getPresion() < 70 && !sensoresEnFuga.contains(m.sector + m.sensor)) {// 70 valor arbitrario para ver si anda JAVA
                filterList.add(m);
                sensoresEnFuga.add(m.sector + m.sensor);
            }
        }
        System.out.println("Mediciones que quedaron despues de filtrar: "+filterList.size());
        return filterList;
    }


    public static class Medicion {

        private String datetime;
        private String sensor;  
        private String sector;
        private int presion;

        public Medicion (){};

        public Medicion (String datetime, String sensor, String sector, int presion) {
            this.datetime = datetime;
            this.sensor = sensor;
            this.sector = sector;
            this.presion = presion;
        }
    
        @JsonProperty("Datetime")
        public void setDatetime(String datetime) {
            this.datetime = datetime;
        }
        
        @JsonProperty("Datetime")
        public String getDatetime() {
            return datetime;
        }
    
        @JsonProperty("Sensor")
        public void setSensor(String sensor) {
            this.sensor = sensor;
        }

        @JsonProperty("Sensor")
        public String getSensor() {
            return sensor;
        }
    
        @JsonProperty("Sector")
        public void setSector(String sector) {
            this.sector = sector;
        }

        @JsonProperty("Sector")
        public String getSector() {
            return sector;
        }
    
        @JsonProperty("Presion")
        public void setPresion(int presion) {
            this.presion = presion;
        }

        @JsonProperty("Presion")
        public int getPresion() {
            return presion;
        }
    }
    
   /*  public static class MyHttpHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            // Lógica para manejar la solicitud GET y devolver la respuesta
            String response = App.getAlertas();
            exchange.sendResponseHeaders(200, response.length());
            OutputStream outputStream = exchange.getResponseBody();
            outputStream.write(response.getBytes());
            outputStream.close();
        }
    } */

    

}


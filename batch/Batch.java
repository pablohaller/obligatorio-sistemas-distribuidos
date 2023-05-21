package batch;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.io.OutputStream;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;
import java.net.InetSocketAddress;

public class Batch {
    public static void main(String[] args) throws Exception {
        int port = 8080;
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        server.createContext("/Batching", new MyHttpHandler());
        server.setExecutor(null); // Utiliza el executor por defecto
        server.start();
        System.out.println("Servidor iniciado en el puerto " + port);
    }



    public static class MyHttpHandler implements HttpHandler {
        @Override
        public void handle(HttpExchange exchange) throws IOException {
            // Lógica para manejar la solicitud GET y devolver la respuesta
            String response = getMediciones();
            exchange.sendResponseHeaders(200, response.length());
            OutputStream outputStream = exchange.getResponseBody();
            outputStream.write(response.getBytes());
            outputStream.close();
        }

        public String getMediciones() {
            String mediciones = "";
            try {
    
                // Get the value of the environment variable
                String mirrorDbServer = System.getenv("MIRROR-DB-SERVER");
    
                // Create a URL object with the endpoint you want to send the request to
                URL url = new URL("http://" + mirrorDbServer + ":8080/UltMediciones/5");
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
                mediciones = response.toString();
    
                // Close the connection
                connection.disconnect();
            } catch (IOException e) {
                e.printStackTrace();
            }
            return mediciones;
        }
    }
}


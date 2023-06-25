# Obligatorio Sistemas Distribuidos

## Manual de opreación

### Requisitos

Para instalar el proyecto hay que simplemente tener instalado Docker Desktop


### Ejecución

En una consola que está parada en el directorio raíz del proyecto, ingresar el comando:

```bash
docker compose up 
```


## Consideraciones

- Para crear un nuevo sector, es necesario crear también una nueva queue para que envíe las mediciones por allí:

Para crear la queue:
```yml
<nombre queue>:
    container_name: <nombre queue>
    image: rabbitmq:3-alpine
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 10s
      timeout: 10s
      retries: 3
    ports:
      - "5672:5672"
    expose:
      - "5672"
```
Para crear el sector:
```yml
# Sector Server
  <nombre sector>:
    container_name: <nombre sector>
    build: ./sector-server
    environment:
      - SECTOR_NAME=<nombre sector>
      - QUEUE_HOST=<nombre queue>
      - QUEUE_NAME=<nombre queue>
      # Coordenadas de ejemplo
        # Latitud y longitud separadas por ,
        # Entre coordenadas la separación es por ;
        # Puedes poner un largo arbitratrio de coordenadas
        # Asegúrate de poner las artistas en orden y que formen un polígono
      - COORDS=-34.91739651002616, -56.16210771871581;-34.91520592550176, -56.15995194170265;-34.915812960803784, -56.15930842617635;-34.91799473178162, -56.161485653707054 
    depends_on:
      <nombre queue>:
        condition: service_healthy
```

- Siempre puedes agregar nuevos sensores a un determinado sector, agregando en el docker compose lo siguiente:

```yml
# Sector Sensor
  <nombre sensor>:
    container_name: <nombre sensor>
    build: ./sensor
    environment:
      - SECTOR_NAME=<sector existente>
      - SENSOR_NAME=<nombre sensor>
      - COORD=<-34.888160, -56.160388>  # Coordenadas de ejemplo
      - MIN_PRESSURE=7.05               # Presion mínima de ejemplo
      - INIT_PRESSURE=13.51             # Presion inicial de ejemplo
    depends_on:
      - <sector existente>  
```

import requests
import json
from datetime import datetime
import time
import random
import os
import pytz
import logging

SectorName = os.environ.get("SECTOR_NAME")
SensorName = os.environ.get("SENSOR_NAME")
MinPressure = os.environ.get("MIN_PRESSURE")
Coord = os.environ.get("COORD")

url = "http://"+SectorName+":8080"  # Reemplaza con la URL correcta del endpoint en Go.
error = random.randint(5, 10)
errorFlag = True
i=0
pressure = 112
# Obtener la zona horaria de Uruguay
timezone = pytz.timezone('America/Montevideo')

sensor = {
    "sensor": SensorName,
    "sector": SectorName,
    "min_pressure": MinPressure,
    "coord": Coord 
}
payloadSensor = json.dumps(sensor)
headersSensor = {
    "Content-Type": "application/json"
    
}

try:
    # Realizar la solicitud PUT al endpoint en Go
    response = requests.put(url + "/Sensor/Suscribe", data=payloadSensor, headers=headersSensor)
    if response.status_code == 200:
        print("Solicitud exitosa")
        logging.info("Solicitud exitosa")
    else:
        print("Error en la solicitud:"+ response.status_code)
        logging.error("Error en la solicitud:" + response.status_code)
except:
    print("An exception occurred")


while True:
    if i == error and errorFlag:
        pressure-=50
        errorFlag = False
    current_datetime = datetime.now(timezone).strftime("%Y-%m-%dT%H:%M:%SZ")
    measurement = {
        "datetime": current_datetime,
        "sensor": SensorName,
        "sector": SectorName,
        "pressure": pressure
    }

    # Convertir el diccionario en una cadena JSON
    payload = json.dumps(measurement)

    # Establecer los encabezados requeridos
    headers = {
        "Content-Type": "application/json"
        
    }
    try:
        # Realizar la solicitud PUT al endpoint en Go
        response = requests.put(url + "/Measurement", data=payload, headers=headers)
        if response.status_code == 200:
            print("Solicitud exitosa")
            logging.info("Solicitud exitosa")
        else:
            print("Error en la solicitud:"+ response.status_code)
            logging.error("Error en la solicitud:" + response.status_code)
    except:
        print("An exception occurred")

    # Esperar 1 segundo antes de realizar la siguiente solicitud
    i+=1
    time.sleep(10)

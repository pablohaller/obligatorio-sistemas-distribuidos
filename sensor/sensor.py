import requests
import json
from datetime import datetime
import time
import random
import os
import pytz
import logging

logging.basicConfig(level=logging.DEBUG)  # Set the log level to DEBUG

SectorName = os.environ.get("SECTOR_NAME")
SensorName = os.environ.get("SENSOR_NAME")
MinPressure = float(os.environ.get("MIN_PRESSURE"))
Coord = os.environ.get("COORD")
InitPressure = float(os.environ.get("INIT_PRESSURE"))

url = "http://"+SectorName+":8080"  # Reemplaza con la URL correcta del endpoint en Go.
error = random.randint(10, 15)
restore = random.randint(25, 35)
errorFlag = True
i=1
pressure = InitPressure
# Obtener la zona horaria de Uruguay
timezone = pytz.timezone('America/Montevideo')

sensor = {
    "sensor": SensorName,
    "sector": SectorName,
    "min_pressure": float(MinPressure),
    "coord": Coord 
}
payloadSensor = json.dumps(sensor)
headersSensor = {
    "Content-Type": "application/json" 
}

# START AddSensor Request #

INIT_BACKOFF = 2
MAX_BACKOFF = 20

retries = 0
while True:
    try:
        # Realizar la solicitud PUT al endpoint en Go
        logging.info("Intentando agregar sensor")
        response = requests.put(url + "/Sensor/Suscribe", data=payloadSensor, headers=headersSensor)
        if response.status_code == 200:
            print("Solicitud de agregar sensor exitosa")
            logging.info("Solicitud de agregar sensor exitosa")
            break
        else:
            print("Error en la solicitud:"+ response.status_code)
            logging.error("Error en la solicitud de agregar sensor:" + response.status_code)
    except requests.exceptions.RequestException as e:
        print("Error en la solicitud:", e)
        logging.error("Error en la solicitud de agregar sensor:" + str(e))
    
    retries += 1
    time.sleep(random.randint(0, min(MAX_BACKOFF, INIT_BACKOFF * pow(2,retries)))) 

# FINISH AddSensor Request #

logging.info(SensorName + " sending measurements since now.")
while True:
    pressure += float(random.randint(-9, 9)) * random.uniform(0.01, 0.1) 
    if (i == error and errorFlag) or pressure <= MinPressure:
        errorFlag = False
    if not errorFlag:
        pressure = float(pressure % MinPressure)

    if (i == restore):
        pressure = InitPressure
        errorFlag = True
        i = 1
        restore = random.randint(18, 30)
        error = random.randint(5, 10)

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
        logging.info("Sending measurement: " + payload)
        if response.status_code == 200:
            print("Solicitud exitosa")
            logging.info("Solicitud exitosa")
        else:
            print("Error en la solicitud:"+ response.status_code)
            logging.error("Error en la solicitud:" + response.status_code)
    except:
        print("An exception occurred")

    # Esperar 10 segundos antes de realizar la siguiente solicitud
    i+=1
    time.sleep(10)

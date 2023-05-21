import requests
import json
from datetime import datetime
import time
import random
import os
import logging

SectorName = os.environ.get("SECTOR_NAME")
SensorName = os.environ.get("SENSOR_NAME")

url = "http://"+SectorName+":8080/Medicion"  # Reemplaza con la URL correcta del endpoint en Go.
error = random.randint(5, 10)
errorFlag = True
i=0
presion = 112
while True:
    if i == error and errorFlag:
        presion-=50
        errorFlag = False

    current_datetime = datetime.now().strftime("%Y-%m-%dT%H:%M:%SZ")
    medicion = {
        "Datetime": current_datetime,
        "Sensor": SensorName,
        "Sector": SectorName,
        "Presion": presion
    }

    # Convertir el diccionario en una cadena JSON
    payload = json.dumps(medicion)

    # Establecer los encabezados requeridos
    headers = {
        "Content-Type": "application/json"
    }


    try:
        # Realizar la solicitud PUT al endpoint en Go
        response = requests.put(url, data=payload, headers=headers)
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
    time.sleep(5)

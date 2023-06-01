FROM openjdk:11

WORKDIR /app

# Copia los archivos fuente y las bibliotecas de Jackson JSON al directorio de trabajo
COPY ./src/App.java /app
COPY jackson-core-2.15.1.jar /app
COPY jackson-databind-2.15.1.jar /app
COPY jackson-annotations-2.15.1.jar /app

# Compila el archivo fuente utilizando las bibliotecas de Jackson JSON
RUN javac -cp jackson-core-2.15.1.jar:jackson-databind-2.15.1.jar:jackson-annotations-2.15.1.jar App.java

# Establece el comando de ejecución
CMD ["java", "-cp", "jackson-core-2.15.1.jar:jackson-databind-2.15.1.jar:jackson-annotations-2.15.1.jar:/app", "App.java"]

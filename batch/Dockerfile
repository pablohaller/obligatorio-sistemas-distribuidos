# Imagen base de Java
FROM openjdk:11

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el archivo .java al contenedor
COPY Batch.java .

# Compilar el archivo .java
RUN javac Batch.java

# Ejecutar el archivo compilado
CMD ["java", "Batch.java"]
# Usa la imagen oficial de Golang como base
FROM golang:1.20-alpine

# Setea el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia el código fuente al contenedor
COPY . .

# Instala las dependencias
RUN go mod tidy

# Exponemos el puerto en el que se va a ejecutar la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["go", "run", "main.go"]
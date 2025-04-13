# Etapa de compilación
FROM golang:1.20-alpine AS build

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar los archivos go.mod y go.sum y descargar las dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente al directorio de trabajo
COPY . .

# Compilar el binario
RUN go build -o main ./cmd/server

# Etapa de producción
FROM alpine:latest

# Crear un directorio en el contenedor para la app
WORKDIR /root/

# Copiar el binario compilado desde la etapa anterior
COPY --from=build /app/main .

# Copiar el archivo .env
COPY --from=build /app/.env .

# Exponer el puerto que usa la app
EXPOSE 3000

# Comando para ejecutar el binario
CMD ["./main"]

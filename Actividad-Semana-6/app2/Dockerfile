# Usar una imagen base oficial de Go
FROM golang:1.16

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos del módulo Go y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto de los archivos del código fuente
COPY . .

# Compilar la aplicación
RUN go build -o main .

# Comando para ejecutar la aplicación
CMD ["./main"]

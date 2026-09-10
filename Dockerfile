FROM golang:1.26.5-alpine

WORKDIR /app

# Copiar go.mod y go.sum para cachear dependencias
COPY go.mod ./
RUN go mod download
# Lo de arriba queda cacheado,ya que no cambiara seguido

COPY . .

# Hacer test sin compilar o buildear el main.
# Agregar al make Script que borre volumenes

RUN go build -o app .

CMD ["./app"]
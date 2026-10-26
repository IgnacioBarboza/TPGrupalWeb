# Etapa de build
FROM golang:1.26.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

# Etapa final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/static ./static
CMD ["./app"]
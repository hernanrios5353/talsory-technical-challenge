# ---------- Build Stage ----------
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiamos go.mod primero para aprovechar el cache de Docker
COPY go.mod ./

RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-go ./cmd/main.go

# ---------- Runtime Stage ----------
FROM alpine:3.20

WORKDIR /app

# Copiamos el binario compilado
COPY --from=builder /app/bin/api-go /app/api-go

# Exponemos el puerto
EXPOSE 8080

# Comando de arranque
CMD ["/app/api-go"]
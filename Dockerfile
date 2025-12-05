# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copiar código
COPY . .

# Compilar
RUN CGO_ENABLED=0 GOOS=linux go build -o puntosgo ./cmd/puntosgo

# Runtime stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/puntosgo .

EXPOSE 3006

CMD ["./puntosgo"]

# Stage 1: Build
FROM golang:1.23 AS builder

# Устанавливаем зависимости для сборки с CGO и librdkafka
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    gcc \
    g++ \
    make \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник с включенным CGO
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./app/main.go

# Stage 2: Run
FROM debian:bookworm-slim

# Устанавливаем runtime-зависимости для librdkafka
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .
COPY config/.env ./config/.env
COPY config/config.yaml ./config/config.yaml
COPY mail/layouts ./mail/layouts

CMD ["./main"]

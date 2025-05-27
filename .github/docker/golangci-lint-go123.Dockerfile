# Stage 1: Prepare modules cache with Go 1.23
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем только файлы с зависимостями для кеша модулей
COPY go.mod go.sum ./

# Загружаем зависимости (кешируем)
RUN go mod download

# Копируем весь код для линтера (только для примера, если нужно)
COPY . .

# Stage 2: golangci-lint based on official image
FROM golangci/golangci-lint:v1.55.2

WORKDIR /app

# Копируем кешированные модули из builder
COPY --from=builder /go/pkg/mod /go/pkg/mod

# Копируем код
COPY --from=builder /app /app

# Устанавливаем переменную окружения для GOMODCACHE, чтобы линтер видел модули
ENV GOMODCACHE=/go/pkg/mod

ENTRYPOINT ["golangci-lint"]
CMD ["run", "./..."]

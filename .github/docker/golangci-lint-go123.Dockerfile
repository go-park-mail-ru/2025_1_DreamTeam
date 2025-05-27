# Stage 1: получить Go 1.23
FROM golang:1.23 as golang

# Stage 2: взять golangci-lint и заменить toolchain
FROM golangci/golangci-lint:v1.55.2

# Копируем Go 1.23 поверх стандартного Go
COPY --from=golang /usr/local/go /usr/local/go

# Обновляем переменные окружения
ENV PATH="/usr/local/go/bin:${PATH}"

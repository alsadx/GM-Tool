FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o sso ./cmd/sso
# Компиляция мигратора
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator
# Скачиваем grpc_health_probe
RUN apk update && apk add wget && \
    wget https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.12/grpc_health_probe-linux-amd64 -O /usr/local/bin/grpc_health_probe && \
    chmod +x /usr/local/bin/grpc_health_probe

# Финальный этап: минимальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем скомпилированные бинарники
COPY --from=builder /app/migrator /migrator

COPY --from=builder /app/sso .

# Копируем прото файлы
COPY protos/ /app/protos/

# Копируем конфигурационный файл
COPY config/local.yaml /app/config/local.yaml

# Копируем миграции
COPY migrations/ /app/migrations/

# Порт, на котором работает SSO-сервис
EXPOSE 8080

# Команда для миграций и запуска сервиса
CMD ["/bin/sh", "-c", "sleep 10 && /migrator && ./sso"]

# RUN go build -o sso ./cmd/sso

# CMD [ "./sso" ]

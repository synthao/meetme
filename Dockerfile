# Этап сборки
FROM golang:1.22-alpine AS builder

# Установка необходимых инструментов для сборки
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов go.mod и go.sum
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование всего исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/app

# Финальный этап
FROM alpine:3.19

# Установка ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

# Создание непривилегированного пользователя
RUN adduser -D -g '' appuser

# Установка рабочей директории
WORKDIR /app

# Копирование исполняемого файла из этапа сборки
COPY --from=builder /app/main .

# Копирование .env файла
# COPY .env .

COPY --from=builder /app/migrations ./migrations/

# Смена владельца файлов
RUN chown -R appuser:appuser /app

# Переключение на непривилегированного пользователя
USER appuser

# Открытие порта
EXPOSE 8080

# Запуск приложения
CMD ["./main"]
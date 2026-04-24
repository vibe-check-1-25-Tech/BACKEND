# Используем официальный образ Go для сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код проекта
COPY . .

# Собираем приложение в один файл под названием "main"
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Финальный легкий образ для запуска
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

# Наш сервер работает на порту 8080
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]
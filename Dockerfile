# Устанавливаем базовый образ с Go для сборки
FROM golang:1.23.1-alpine AS build

LABEL name="test_avtito"
# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY . .

# Загружаем зависимости
RUN go mod download

# Собираем проект, указывая путь к main.go в папке server
RUN go build -o main ./cmd/server/main.go

# Минимальный образ для выполнения
FROM alpine:latest

# Создаём директорию для приложения
WORKDIR /app

# Копируем бинарный файл из стадии сборки
COPY --from=build /app/main .

# Экспортируем порт, на котором будет запущен сервер
EXPOSE 8080

# Указываем команду для запуска сервера
ENTRYPOINT ["./main"]
# Этап сборки
FROM golang:1.20.1-alpine3.16 as base



# Установка рабочей директории
WORKDIR /web

# Копируем все файлы в контейнер
COPY go.mod go.sum ./

# Строим проект
RUN go mod download

COPY . .

RUN mkdir -p /web/logging
RUN mkdir -p /web/pkg/logger

RUN go build -o web ./cmd/main.go


EXPOSE 8080 

# Указываем команду для запуска контейнера
CMD ["./web"]

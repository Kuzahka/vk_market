# Используем официальный образ Go для сборки
FROM golang:1.24-alpine AS builder


WORKDIR /app


COPY go.mod .
COPY go.sum .


RUN go mod download


COPY . .


RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/main ./cmd/api


FROM alpine:latest


WORKDIR /app


COPY --from=builder /app/main .

COPY .env .


EXPOSE 8080


CMD ["./main"]

FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main ./
COPY ./config ./config
COPY ./family.json ./family.json

EXPOSE 8081

CMD ["./main"]
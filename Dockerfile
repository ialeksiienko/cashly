FROM golang:1.23-alpine AS builder

WORKDIR /cashly

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /cashly

COPY --from=builder /cashly/main ./
COPY ./config ./config
COPY ./family.json ./family.json

EXPOSE 8081

CMD ["./main"]

# FROM golang:1.23.0 as builder
# WORKDIR /go/src/app

# ENV GO111MODULE=on

# RUN go install github.com/cespare/reflex@latest

# RUN go get github.com/google/uuid
# RUN go get github.com/rabbitmq/amqp091-go
# RUN go get github.com/ilyakaznacheev/cleanenv
# RUN go get gopkg.in/telebot.v3

# COPY go.mod .
# COPY go.sum .

# RUN go mod tidy
# RUN go mod download

# COPY cashly/. cashly/.

# RUN go build -o ./run ./cashly/cmd/.

# FROM alpine:latest
# WORKDIR /root/

# COPY --from=builder /go/src/app/run .

# CMD ["./run"]

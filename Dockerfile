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
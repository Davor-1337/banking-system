FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o db-service ./db-service

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/db-service .

EXPOSE 50051

CMD ["./db-service"]
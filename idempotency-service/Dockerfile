
FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN go build -o idempotency-service ./idempotency-service

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/idempotency-service .

EXPOSE 50052

CMD ["./idempotency-service"]
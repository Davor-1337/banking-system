
FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

WORKDIR /app/web-server

RUN go build -o web-server .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/web-server/web-server .

EXPOSE 8080

CMD ["./web-server"]
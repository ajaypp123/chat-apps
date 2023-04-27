FROM golang:1.16.5-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o client ./cmd/client/client.go

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/client .
CMD ["./client"]
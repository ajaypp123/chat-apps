FROM golang:1.16.5-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server ./cmd/server/chat-apps.go

FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 50051
CMD ["./server"]


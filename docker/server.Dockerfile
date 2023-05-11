# Use the official Golang image as the base image
FROM golang:1.18 AS builder
# Set the working directory inside the container
WORKDIR /chat-server
# Copy the source code into the container
COPY . .

# Build the Go binary with the necessary flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/chat-server -mod vendor ./cmd/server/chat-apps.go

# Use a small Alpine-based image as the base image for the final container
FROM alpine:latest AS deploy
# Set the working directory inside the container
WORKDIR /chat-server
RUN mkdir ./configs/

# Copy the built binary from the builder container
COPY --from=builder /chat-server/build/chat-server .
COPY --from=builder /chat-server/configs/config.json ./configs/

# Expose the port that the server listens on
EXPOSE 8080

# Start the server when the container starts
CMD ["./chat-server -grpc 50050 -http 8080"]

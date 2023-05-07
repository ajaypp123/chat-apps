# chat-apps
This is a distributed chat application built using Golang and gRPC. The application allows multiple clients to connect to multiple servers and exchange messages in real-time.

The servers use Redis to store client information, allowing clients to connect to any server and maintain their session. The application also uses Kafka for Pub/Sub messaging, allowing servers to communicate with each other and route messages to the appropriate client.

The project includes a proto file defining the gRPC service, as well as server and client implementations. The client includes a simple command-line interface for sending and receiving messages.

# Features
- One on one chat
- Chat history
- Ability to add more features as needed

# Technologies Used
- gRPC for communication
- Golang for building the application
- Redis for storing the gRPC session
- Distributed client-server architecture for enabling communication

# Target Audience
- Individuals who want to experiment with building chat applications
- Companies who want to use a basic chat application for their internal communication needs
- Anyone who wants to learn about gRPC and distributed systems

# Instructions for Running the Application
- Start the server by running the following command in the terminal:
```
go run cmd/server/chat-app.go
```

- Start multiple instances of the client for each user. Run the following command for each instance:

```
go run cmd/client/client.go
```
- Once the client is started, enter the username when prompted. After registration, enter the sender's name and message to send the message. The client will also receive incoming messages.

# Docker build
```
docker-compose up --build
```

Feel free to experiment and add more features to the application as needed.


# Start and use tools for startup guide

1. start service
```
cd chat-apps/
go mod tidy
go mod vendor
```

1. start server and create multiple user
```
go run cmd/server/chat-apps.go

# create users
curl -XPOST localhost:8080/v1/chat-apps/users?username=A

# check user detail as per need
curl -X POST localhost:8080/v1/chat-apps/users -d '{ "username":"ajuser","name":"bar", "phone": "8098080"}'

{"req_id":"","status":"success","data":{"username":"ajuser","name":"bar","phone":"8098080","secret":"cab3b16f-d9e5-42fb-965d-2321893613de"},"code":200}

curl -XGET localhost:8080/v1/chat-apps/users?username=ajuser

{"req_id":"","status":"success","data":{"username":"ajuser","name":"bar","phone":"8098080","secret":"cab3b16f-d9e5-42fb-965d-2321893613de"},"code":200}
```

# TODO: There is room to grow for app
1. queue for dead message not verified
2. group message
3. unread message utility

# chat-apps
This is a distributed chat application built using Golang and gRPC. The application allows multiple clients to connect to multiple servers and exchange messages in real-time.

The servers use Redis to store client information, allowing clients to connect to any server and maintain their session. The application also uses Kafka for Pub/Sub messaging, allowing servers to communicate with each other and route messages to the appropriate client.

The project includes a proto file defining the gRPC service, as well as server and client implementations. The client includes a simple command-line interface for sending and receiving messages.

## Features
- One on one chat
- Chat history
- Ability to add more features as needed

## Technologies Used
- gRPC for communication
- Golang for building the application
- Redis for storing the gRPC session
- Distributed client-server architecture for enabling communication

## Target Audience
- Individuals who want to experiment with building chat applications
- Companies who want to use a basic chat application for their internal communication needs
- Anyone who wants to learn about gRPC and distributed systems

## Instructions for Running the Application in Docker
- It starts 2 server, and then we can run cli for each client

1. Start the server by running the following command in the terminal:
```
cd chat-apps/docker
docker-compose up -d
```

1. Register users by sending POST requests to the following endpoint:
```shell
curl -X POST localhost:3000/v1/chat-apps/users -d '{ "username":"user1","name":"bar", "phone": "8098080"}'

curl -X POST localhost:3000/v1/chat-apps/users -d '{ "username":"user2","name":"bar", "phone": "8098080"}'

curl -X POST localhost:3000/v1/chat-apps/users -d '{ "username":"user3","name":"bar", "phone": "8098080"}'

curl -XGET localhost:3000/v1/chat-apps/users?username=user1
```

1. Start the client by running the following command in the terminal:
```
go run cmd/client/client.go -grpc :50055 -http :3000
```

1. Once the client is started, enter the username when prompted. After registration, enter the sender's name and message to send the message. The client will also receive incoming messages.

## Instructions for Running the Application without Docker

- Pre-requisite
1. Install redis
2. Update connection detail in config.json

- Start application
1. start service
```
cd chat-apps/
go mod tidy
go mod vendor
```

1. start server and create multiple user
```
<<<<<<< HEAD
go run cmd/server/chat-apps.go -grpc :50050 -http :3000
=======
go run cmd/server/chat-apps.go

# create users
curl -XPOST localhost:8080/v1/chat-apps/users?username=A

# check user detail as per need
curl -X POST localhost:8080/v1/chat-apps/users -d '{ "username":"ajuser","name":"bar", "phone": "8098080"}'

{"req_id":"","status":"success","data":{"username":"ajuser","name":"bar","phone":"8098080","secret":"cab3b16f-d9e5-42fb-965d-2321893613de"},"code":200}

curl -XGET localhost:8080/v1/chat-apps/users?username=ajuser

{"req_id":"","status":"success","data":{"username":"ajuser","name":"bar","phone":"8098080","secret":"cab3b16f-d9e5-42fb-965d-2321893613de"},"code":200}
>>>>>>> origin/main
```

1. register users
```shell
curl -X POST localhost:3000/v1/chat-apps/users -d '{ "username":"user1","name":"bar", "phone": "8098080"}'

curl -X POST localhost:3000/v1/chat-apps/users -d '{ "username":"user2","name":"bar", "phone": "8098080"}'

curl -X POST localhost:3000/v1/chat-apps/users -d '{ "username":"user3","name":"bar", "phone": "8098080"}'

curl -XGET localhost:3000/v1/chat-apps/users?username=user1
```

1. Start client
```
go run cmd/client/client.go -grpc :50055 -http :3000
```

1. Once the client is started, enter the username when prompted. After registration, enter the sender's name and message to send the message. The client will also receive incoming messages.


# TODO: There is room to grow for app
1. queue for dead message not verified
2. group message
3. unread message utility

username should be metadata so that during failure, we will be able to close connection

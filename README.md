# chat-apps
chat-apps is a basic chat application that enables one on one communication between users. It uses gRPC for communication and is built with Golang. The application stores the gRPC session in Redis and can be run in a distributed client-server architecture.

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
go run cmd/server/chat-apps.go -grpc :50050 -http :3000
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

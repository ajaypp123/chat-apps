# chat-apps
chat-apps is a basic chat application that enables one on one communication between users. It uses gRPC for communication and is built with Golang. The application stores the gRPC session in Redis and can be run in a distributed client-server architecture.

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

#go run cmd/server/chat-apps.go -grpc :50051 -http :8081
#go run cmd/server/chat-apps.go -grpc :50052 -http :8082

#go run cmd/client/client.go -grpc :50055 -http :9000

# kill -9 $(lsof -t -i :50052)

curl -X POST localhost:9000/v1/chatapp/users -d '{ "username":"user1","name":"bar", "phone": "8098080"}'

curl -X POST localhost:9000/v1/chatapp/users -d '{ "username":"user2","name":"bar", "phone": "8098080"}'

curl -X POST localhost:9000/v1/chatapp/users -d '{ "username":"user3","name":"bar", "phone": "8098080"}'

curl -XGET localhost:9000/v1/chatapp/users?username=user1
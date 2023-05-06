package streamer

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisStreamingService struct {
	client *redis.Client
	pubSub *redis.PubSub
}

var ctx = context.Background()

// NewRedisStreamingService creates a new RedisKVStore instance
func NewRedisStreamingService(addr string, password string, db int) (StreamingServiceI, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Check if Redis is connected
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	cli := &RedisStreamingService{client: client}
	if defaultClient == nil {
		defaultClient = cli
	}

	return cli, nil
}

func (r *RedisStreamingService) PublishMessage(topic, partition string, message string) error {
	return r.client.Publish(ctx, topic+partition, message).Err()
}

func (r *RedisStreamingService) RegisterSubscriber(sub Subscriber) {
	if subscribers == nil {
		subscribers = []Subscriber{}
	}
	subscribers = append(subscribers, sub)
}

func (r *RedisStreamingService) StartListening(topic, partition string) error {
	// Create subscriber
	r.pubSub = r.client.Subscribe(context.Background(), topic+partition)

	// Subscribe to channel
	ch := r.pubSub.Channel()

	go func() {
		// Consume messages from channel
		for msg := range ch {
			//fmt.Println("Received message:", msg.Payload)
			for _, sub := range subscribers {
				sub.ReceiveMsg(msg.Payload)
			}
		}
	}()
	return nil
}

func (r *RedisStreamingService) StopListening() {
	defer r.pubSub.Close()
}

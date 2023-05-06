package streamer

import (
	"context"
	"errors"
)

var subscribers []Subscriber

type Subscriber interface {
	ReceiveMsg(ctx context.Context, msg interface{})
}

type StreamingServiceI interface {
	PublishMessage(topic, partition string, message string) error
}

var defaultClient StreamingServiceI

func GetStreamingService() (StreamingServiceI, error) {
	if defaultClient != nil {
		return defaultClient, nil
	}
	return nil, errors.New("KVStore is not initialised")
}

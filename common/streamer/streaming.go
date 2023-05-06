package streamer

import (
	"errors"
)

var subscribers []Subscriber

type Subscriber interface {
	ReceiveMsg(msg interface{})
}

type StreamingServiceI interface {
	RegisterSubscriber(sub Subscriber)
	PublishMessage(topic, partition string, message string) error
	StartListening(topic, partition string) error
	StopListening()
}

var defaultClient StreamingServiceI

func GetStreamingService() (StreamingServiceI, error) {
	if defaultClient != nil {
		return defaultClient, nil
	}
	return nil, errors.New("KVStore is not initialised")
}

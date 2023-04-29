package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

func GetTimeNow() string {
	return cast.ToString(time.Now().UTC().UnixMilli())
}

func GetUUID() string {
	return uuid.New().String()
}

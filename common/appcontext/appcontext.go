package appcontext

import (
	"context"
	"encoding/json"
)

type AppContext struct {
	context.Context
	kv map[string]interface{}
}

func GetDefaultContext() *AppContext {
	ctx := &AppContext{
		Context: context.Background(),
		kv:      make(map[string]interface{}),
	}
	return ctx
}

// Method to add key-value pairs to AppContext
func (c *AppContext) AddValue(key string, value interface{}) {
	if c.kv == nil {
		c.kv = make(map[string]interface{})
	}
	c.kv[key] = value
}

func (c *AppContext) GetValue(key string) interface{} {
	data, ok := c.kv[key]
	if !ok {
		return nil
	}
	return data
}

// MarshalJSON implements json.Marshaler interface to print key-value pairs as JSON
func (c *AppContext) marshalJSON() ([]byte, error) {
	data := make(map[string]interface{})
	for k, v := range c.kv {
		data[k] = v
	}
	return json.Marshal(data)
}

func (c *AppContext) String() string {
	data, _ := c.marshalJSON()
	return string(data)
}

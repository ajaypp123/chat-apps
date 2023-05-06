package appcontext

import (
	"context"
	"encoding/json"
)

var DefaultContext *AppContext

type AppContext struct {
	context.Context
	kv map[string]interface{}
}

func GetNewContext() *AppContext {
	ctx := &AppContext{
		Context: context.Background(),
		kv:      make(map[string]interface{}),
	}
	return ctx
}

func GetDefaultContext(data map[string]string) *AppContext {
	if DefaultContext != nil {
		return DefaultContext
	}

	ctx := GetNewContext()
	for key, val := range data {
		ctx.AddValue(key, val)
	}

	DefaultContext = ctx
	return ctx
}

func (c *AppContext) DeepCopy() *AppContext {
	// Create a new context that is a copy of the original context
	newCtx := GetNewContext()

	// Create a new key-value map that is a copy of the original map
	for k, v := range c.kv {
		newCtx.kv[k] = v
	}

	// Return a new AppContext with the new context and key-value map
	return newCtx
}

// AddValue Method to add key-value pairs to AppContext
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

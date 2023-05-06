package common

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Config struct{}

func ConfigService() *Config {
	return &Config{}
}

func (c *Config) Init() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	return err
}

func (c *Config) GetValue(key string) string {
	return cast.ToString(viper.Get(key))
}

func (c *Config) GetNodeId() string {
	return cast.ToString(viper.Get("server_id"))
}

func (c *Config) GetStreamDetail() string {
	return cast.ToString(viper.Get("server_id"))
}

package redis_config

import (
	"github.com/nihal-ramaswamy/GoChat/internal/constants"
	"github.com/nihal-ramaswamy/GoChat/internal/utils"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisConfig(options ...func(*RedisConfig)) *RedisConfig {
	config := &RedisConfig{}
	for _, option := range options {
		option(config)
	}
	return config
}

func WithAddrRedis(addr string) func(*RedisConfig) {
	return func(c *RedisConfig) {
		c.Addr = addr
	}
}

func WithPasswordRedis(password string) func(*RedisConfig) {
	return func(c *RedisConfig) {
		c.Password = password
	}
}

func WithDBRedis(db int) func(*RedisConfig) {
	return func(c *RedisConfig) {
		c.DB = db
	}
}

func DefaultRedisAuthConfig() *RedisConfig {
	return NewRedisConfig(
		WithAddrRedis(utils.GetDotEnvVariable(constants.REDIS_HOST)+":"+utils.GetDotEnvVariable(constants.REDIS_PORT)),
		WithPasswordRedis(utils.GetDotEnvVariable(constants.REDIS_PASSWORD)),
		WithDBRedis(0),
	)
}

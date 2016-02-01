package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jgautheron/exago-service/config"
)

var (
	pool *redis.Pool
)

func init() {
	pool = newPool(config.Get("RedisHost") + ":6379")
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			var err error
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

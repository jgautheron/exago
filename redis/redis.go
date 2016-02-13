package redis

import (
	"time"

	"github.com/exago/svc/config"
	"github.com/garyburd/redigo/redis"
)

var (
	pool *redis.Pool
)

func SetUp() {
	pool = newPool(config.Get("RedisHost") + ":6379")
}

func GetConn() redis.Conn {
	return pool.Get()
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			var err error

			c, err := redis.Dial(
				"tcp",
				server,
				redis.DialConnectTimeout(2*time.Second),
			)
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

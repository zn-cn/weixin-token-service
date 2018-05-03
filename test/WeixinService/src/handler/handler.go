package handler

import (
	"github.com/garyburd/redigo/redis"
)

type (
	// Handler handler
	Handler struct {
		RedisConn *redis.Pool
	}
)

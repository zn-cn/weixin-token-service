package daemon

import (
	"conf"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

// Handler handler
type Handler struct {
	RedisConn *redis.Conn
}

// DBInit redis初始化
func DBInit() (*Handler, error) {
	// connection prepare
	// redis prepare
	redisConn, redisErr := redis.Dial("tcp",
		fmt.Sprintf("%s:%s", conf.Conf.Redis.Host, conf.Conf.Redis.Port))
	if redisErr != nil {
		log.Fatalln("redis connection failed")
	}

	// init handler
	h := &Handler{RedisConn: &redisConn}

	return h, nil

}

package daemon

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

// Handler handler
type Handler struct {
	RedisConn [3](*redis.Conn)
}

// DBInit redis初始化
func DBInit() (*Handler, error) {
	// connection prepare
	// redis prepare
	var redisConn [3](*redis.Conn)
	temp1, redisErr := redis.Dial("tcp",
		fmt.Sprintf("redis1:6379"))
	if redisErr != nil {
		log.Fatalln("redis2 connection failed")
	}
	redisConn[0] = &temp1

	temp2, redisErr := redis.Dial("tcp",
		fmt.Sprintf("redis2:6379"))
	if redisErr != nil {
		log.Fatalln("redis2 connection failed")
	}
	redisConn[0] = &temp2

	temp3, redisErr := redis.Dial("tcp",
		fmt.Sprintf("redis3:6379"))
	if redisErr != nil {
		log.Fatalln("redis3 connection failed")
	}
	redisConn[2] = &temp3
	// init handler
	h := &Handler{RedisConn: redisConn}

	return h, nil

}

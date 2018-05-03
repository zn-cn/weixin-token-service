package restful

import (
	"handler"
	"util"

	"github.com/garyburd/redigo/redis"
)

var dbLog = util.GetLogger("/app/log/restful/db.txt", "[DEBUG]")

// DbInit 初始化redis
func DbInit() (*handler.Handler, error) {
	// connection prepare
	// redis prepare
	// 建立连接池
	redisClient := &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:   10,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "redis:6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}

	// redisCon, redisErr := redis.Dial("tcp",
	// 	fmt.Sprintf("%s:%s", conf.Conf.Redis.Host, conf.Conf.Redis.Port))
	// if redisErr != nil {
	// 	dbLog.Println("redis connection failed")
	// }

	// init handler
	h := &handler.Handler{RedisConn: redisClient}

	return h, nil

}

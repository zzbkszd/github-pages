package redis_client

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)

var (
	Conn redis.Conn
	err error
)

func init()  {
	Conn, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
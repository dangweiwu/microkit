package main

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/dangweiwu/microkit/db/redisx"
	"log"
)

func main() {
	config := redisx.Config{}

	redisServer, err := miniredis.Run()

	if err != nil {
		log.Panicf("redisServer创建失败:%v\n", err)
	}

	defer redisServer.Close()
	config.Addr = redisServer.Addr()

	redisClient, err := redisx.NewClient(config)
	if err != nil {
		log.Panicf("redisClient创建失败:%v\n", err)
	}

	redisClient.Set(context.Background(), "hello", "world", 0)
	log.Println(redisClient.Get(context.Background(), "hello").Val())
}

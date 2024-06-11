package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(context.Background(), "192.168.100.1", 50, 0).Err()
	if err != nil {
		fmt.Println(err.Error())
	}

	val, err := rdb.Get(context.Background(), "192.168.100.1").Result()
	if err != nil {
		fmt.Println(err.Error())
	}

	converted, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println(err.Error())
	}
	if converted == 50 {
		fmt.Println(converted)
	}
}

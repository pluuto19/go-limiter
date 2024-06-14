package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	//
	//err := rdb.Set(context.Background(), "a", 1, 0).Err()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//val, err := rdb.Get(context.Background(), "a").Result()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(val)
	//err1 := rdb.Set(context.Background(), "a", 0, 0).Err()
	//if err1 != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//val1, err1 := rdb.Get(context.Background(), "a").Result()
	//if err1 != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//converted, err2 := strconv.Atoi(val1)
	//if err2 != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(converted)
	//now := time.Now()
	//
	//// Format the time according to RFC 1123
	//formattedDate := now.Format(time.RFC1123)
	//
	//// Construct the HTTP response string
	//httpResponse := fmt.Sprintf("HTTP/1.1 429 Too Many Requests\r\nConnection: close\r\nDate: %s\r\n", formattedDate)
	//
	//// Print the HTTP response
	//fmt.Println(httpResponse)

	rdb.ZAdd(context.Background(), "192.168.100.1", redis.Z{
		Score:  float64(time.Now().UnixNano()),
		Member: time.Now().UnixNano(),
	})

	//str := rdb.ZRange(context.Background(), "192.168.100.1", 0, -1)
	//fmt.Println(str.Val())

	str, err := rdb.ZRevRange(context.Background(), "192.168.100.1", 0, -1).Result()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("might be empty")
	}

	fmt.Println(str[0])
	val, err := strconv.Atoi(str[0])
	if err != nil {
		return
	}
	fmt.Println(val)
	fmt.Println(time.Now().UnixNano())
	fmt.Println(time.Now().UnixNano() - int64(val))

	timeBefore := time.Now().UnixNano()
	time.Sleep(3 * 1000000000 * time.Nanosecond)
	timeAfter := time.Now().UnixNano()
	fmt.Println(timeAfter - timeBefore)

}

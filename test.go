package main

import (
	"fmt"
	"time"
)

func main() {
	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
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
	now := time.Now()

	// Format the time according to RFC 1123
	formattedDate := now.Format(time.RFC1123)

	// Construct the HTTP response string
	httpResponse := fmt.Sprintf("HTTP/1.1 429 Too Many Requests\r\nConnection: close\r\nDate: %s\r\n", formattedDate)

	// Print the HTTP response
	fmt.Println(httpResponse)
}

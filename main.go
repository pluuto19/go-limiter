package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
)

var ctx = context.Background()

func main() {
	cache := initializeCacheStore("localhost", "6379", "")

	serverSpec, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	welcSock, err := net.ListenTCP("tcp", serverSpec)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		clientConnSock, err := welcSock.Accept() // client connects to the rate limiter
		if err != nil {
			fmt.Println(err)
			continue
		}
		go manageConnection(clientConnSock, cache) // fetch tokens from the cache store and decide whether to rate limit
	}
}

func initializeCacheStore(addr string, port string, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: pass,
		DB:       0,
	})
	return rdb
}

func manageConnection(clientConnSock net.Conn, cache *redis.Client) {
	forwardRequest := checkValidity(clientConnSock, cache)
	if forwardRequest { // if passed then set correct headers when response comes back from backend

	} else { // if rate limited then return correct HTTP respone

	}
	// and update the tokens in cache
}

func checkValidity(clientConnSock net.Conn, cache *redis.Client) bool {
	// fetch from store
	tokensLeft, err := cache.Get(ctx, clientConnSock.RemoteAddr().String()).Result()
	if err != nil {
		panic(err) // no key
	}
	fmt.Println("key", tokensLeft)

	// run the Token Bucket/Sliding Window Log algorithm to decide whether to rate limit or let it pass
}
func updateCache(clientConnSock net.Conn, cache *redis.Client) {

}

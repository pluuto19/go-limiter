package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

var ctx = context.Background()

const tokenSize = 10
const refillRate = 3

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
	forwardRequest, tokenRemaining := checkValidity(clientConnSock, cache)
	if forwardRequest { // if passed then set correct headers when response comes back from backen
		httpresponse := fmt.Sprintf("HTTP/1.1 429 Too Many Requests\r\nConnection: close\r\nDate: %s\r\nServer: GoLang/1.22.2(Alpine)\r\nX-Ratelimit-Remaining: %s\r\nX-Ratelimit-Limit: %s\r\n\r\n", time.Now().Format(time.RFC1123), strconv.Itoa(tokenRemaining), strconv.Itoa(tokenSize))
	} else { // if rate limited then return correct HTTP response
		httpresponse := fmt.Sprintf("HTTP/1.1 429 Too Many Requests\r\nConnection: close\r\nDate: %s\r\nServer: GoLang/1.22.2(Alpine)\r\nX-Ratelimit-Retry-After: %s", time.Now().Format(time.RFC1123), "")

	}
}

func checkValidity(clientConnSock net.Conn, cache *redis.Client) (bool, int) {
	tokensLeft, err := cache.Get(ctx, clientConnSock.RemoteAddr().String()).Result() // fetch from cache
	if err != nil {
		// set the key (IP addr) to value(tokens-1) in cache
		cache.Set(ctx, clientConnSock.RemoteAddr().String(), tokenSize-1, refillRate*time.Second)
		fmt.Println(err.Error())
		return true, tokenSize - 1
	}
	converted, err := strconv.Atoi(tokensLeft)
	if converted == 0 {
		return false, 0
	} else {
		cache.Decr(ctx, clientConnSock.RemoteAddr().String()) // update the tokens in cache
		return true, converted - 1
	}
}

// run the Token Bucket/Sliding Window Log algorithm to decide whether to rate limit or let it pass

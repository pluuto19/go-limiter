package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

const tokenSize = 10
const refillRate = 3

func checkValidityTokenBucket(clientConnSock net.Conn, cache *redis.Client) (bool, int) {
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

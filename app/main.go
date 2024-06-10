package main

import (
	"fmt"
	"net"
)

func main() {
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
		go fetchFromCache(clientConnSock) // fetch tokens from the cache store and decide whether to rate limit
	}
}

func fetchFromCache(clientConnSock net.Conn) {
	// fetch from store
	// run the Token Bucket/Sliding Window Log algorithm to decide whether to rate limit or let it pass
	// if rate limited then return correct HTTP respone
	// if passed then also set correct headers when response comes back from backend
	// and update the tokens in cache
}

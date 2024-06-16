package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

const bufSize = 1024
const loadBalAddr = ":"

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
	forwardRequest, remLim := checkValidityTokenBucket(clientConnSock, cache)
	if forwardRequest { // if passed then set correct headers when response comes back from backend
		clientRecvBuffer := make([]byte, bufSize)
		clientRecvBufLen, err := clientConnSock.Read(clientRecvBuffer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		resolvedLoadBalAddr, err := net.ResolveTCPAddr("tcp4", loadBalAddr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		loadBalSock, err := net.DialTCP("tcp4", nil, resolvedLoadBalAddr)
		_, err1 := loadBalSock.Write(clientRecvBuffer[0:clientRecvBufLen])
		if err1 != nil {
			fmt.Println(err1.Error())
			return
		}

		loadBalRecvBuffer := make([]byte, bufSize)
		loadBalRecvBufLen, err := loadBalSock.Read(loadBalRecvBuffer)
		closeerr := loadBalSock.Close()
		if closeerr != nil {
			fmt.Println(closeerr.Error())
			return
		}

		fmt.Println(loadBalRecvBuffer[0:loadBalRecvBufLen])
		fmt.Println(loadBalRecvBufLen)
		fmt.Println(string(loadBalRecvBuffer))

		attachedHeadersBuf := attachExtraHeaders(loadBalRecvBuffer, loadBalRecvBufLen, remLim, tokenSize)

		// add the correct headers to the message and send it into the client-socket
		_, writeerr := clientConnSock.Write(attachedHeadersBuf[0:])
		if writeerr != nil {
			fmt.Println(writeerr.Error())
			return
		}

		clientcloseerr := clientConnSock.Close()
		if clientcloseerr != nil {
			fmt.Println(clientcloseerr.Error())
			return
		}

	} else { // if rate limited then return correct HTTP response
		httpresponse := fmt.Sprintf("HTTP/1.1 429 Too Many Requests\r\nConnection: close\r\nDate: %s\r\nServer: GoLang/1.22.2(Alpine)\r\nX-Ratelimit-Retry-After: %s", time.Now().Format(time.RFC1123), "") // mostrecenttimestamp + refillrate(or refillafter) - currentrequesttimestamp
		_, err := clientConnSock.Write([]byte(httpresponse))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}
}

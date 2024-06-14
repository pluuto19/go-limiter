package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
)

const refillAfter = 3000000000 // 3 seconds in UnixNano()
const requestsAllowed = 5

func checkValiditySWindowLog(clientConnSock net.Conn, cache *redis.Client, requestArrivalTime int64) (bool, int) {
	// fetch log from cache
	resultLog, err := cache.ZRevRange(ctx, clientConnSock.RemoteAddr().String(), 0, -1).Result()
	if err != nil {
		fmt.Println(err.Error())
		return false, 0
	}
	if len(resultLog) >= 1 {
		converted, err := strconv.ParseInt(resultLog[0], 10, 64) // convert the most recent request's timestamp to int64
		if err != nil {
			fmt.Println(err.Error())
			return false, 0
		}

		if requestArrivalTime-converted >= refillAfter { // request arriving after refillRate time has elapsed
			cache.Del(ctx, clientConnSock.RemoteAddr().String())
			cache.ZAdd(ctx, clientConnSock.RemoteAddr().String(), redis.Z{
				Score:  float64(requestArrivalTime),
				Member: requestArrivalTime,
			})
			return true, requestsAllowed - 1
		} else { // request arriving within the time elapsed
			if len(resultLog) < requestsAllowed {
				cache.ZAdd(ctx, clientConnSock.RemoteAddr().String(), redis.Z{
					Score:  float64(requestArrivalTime),
					Member: requestArrivalTime,
				})
				return true, requestsAllowed - len(resultLog) - 1
			} else {
				return false, 0
			}
		}
	} else {
		cache.ZAdd(ctx, clientConnSock.RemoteAddr().String(), redis.Z{
			Score:  float64(requestArrivalTime),
			Member: requestArrivalTime,
		})
		return true, requestsAllowed - 1
	}
}

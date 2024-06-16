package main

import (
	"bytes"
	"fmt"
	"strconv"
)

// httpresponse := fmt.Sprintf("HTTP/1.1 200 OK\r\nConnection: close\r\nDate: %s\r\nServer: GoLang/1.22.2(Alpine)\r\nX-Ratelimit-Remaining: %s\r\nX-Ratelimit-Limit: %s\r\n\r\n", time.Now().Format(time.RFC1123), strconv.Itoa(tokenRemaining), strconv.Itoa(tokenSize))

func attachExtraHeaders(htmlRecvBuffer []byte, bufSize int, remRateLim int, rateLim int) []byte {
	breakIndex := bytes.Index(htmlRecvBuffer, []byte("\r\n\r\n"))
	if breakIndex == -1 {
		fmt.Println("abc")
		return []byte("")
	}
	headers := htmlRecvBuffer[:breakIndex]
	body := htmlRecvBuffer[breakIndex+4:]
	//bufSize - len(headers)
	newHeaders := fmt.Sprintf("X-Ratelimit-Remaining: %s\r\nX-Ratelimit-Limit: %s\r\n", strconv.Itoa(remRateLim), strconv.Itoa(rateLim))
	modifiedHeaders := string(headers) + "\r\n" + newHeaders

	return []byte(modifiedHeaders + "\r\n" + string(body))
}

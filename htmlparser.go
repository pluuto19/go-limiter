package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func attachExtraHeaders(htmlRecvBuffer []byte, bufSize int, remRateLim int, rateLim int) []byte {
	breakIndex := bytes.Index(htmlRecvBuffer, []byte("\r\n\r\n"))
	if breakIndex == -1 {
		fmt.Println("Incorrect HTML Response Format returned by Web Server...")
		return []byte("")
	}
	headers := htmlRecvBuffer[:breakIndex]
	body := htmlRecvBuffer[breakIndex+4:]
	newHeaders := fmt.Sprintf("X-Ratelimit-Remaining: %s\r\nX-Ratelimit-Limit: %s\r\n", strconv.Itoa(remRateLim), strconv.Itoa(rateLim))
	modifiedHeaders := string(headers) + "\r\n" + newHeaders

	return []byte(modifiedHeaders + "\r\n" + string(body))
}

package main

// httpresponse := fmt.Sprintf("HTTP/1.1 200 OK\r\nConnection: close\r\nDate: %s\r\nServer: GoLang/1.22.2(Alpine)\r\nX-Ratelimit-Remaining: %s\r\nX-Ratelimit-Limit: %s\r\n\r\n", time.Now().Format(time.RFC1123), strconv.Itoa(tokenRemaining), strconv.Itoa(tokenSize))

func attachExtraHeaders(htmlRecvBuffer []byte, bufSize int, remRateLim int, rateLim int) ([]byte, int) {

}

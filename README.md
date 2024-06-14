An HTTP rate limiter as a separate service/middleware and implements the following algorithms:

- Token Bucket
- Sliding Window Log

The service uses Redis as the in-memory cache and establishes a connection with  a load balancer
to further propagate the HTTP request.
It also let's the client know about it's rate limiting parameters on each request.
If the client gets rate limited, it simply returns an HTTP 429.
Further improvements could include sending the request to a messaging service.
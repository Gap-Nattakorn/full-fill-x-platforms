package redis

type Client struct {
	URL string
}

func NewClient(redisURL string) *Client {
	return &Client{URL: redisURL}
}

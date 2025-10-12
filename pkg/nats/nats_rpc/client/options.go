package client

import "time"

// Option -.
type Option func(*Client)

// Timeout -.
func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

package client

import "time"

type Option func(*Client)

func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func ConnWaitTime(timeout time.Duration) Option {
	return func(c *Client) {
		c.conn.WaitTime = timeout
	}
}

func ConnAttempts(attempts int) Option {
	return func(c *Client) {
		c.conn.Attempts = attempts
	}
}

package natsrpc

import "github.com/nats-io/nats.go"

// HeaderCarrier adapts nats.Header to the OpenTelemetry propagation.TextMapCarrier
// interface so trace context can be injected into / extracted from NATS message headers.
type HeaderCarrier nats.Header

// Get returns the value associated with the passed key.
func (c HeaderCarrier) Get(key string) string {
	values, ok := c[key]
	if !ok || len(values) == 0 {
		return ""
	}

	return values[0]
}

// Set stores the key-value pair.
func (c HeaderCarrier) Set(key, value string) {
	c[key] = []string{value}
}

// Keys lists the keys stored in this carrier.
func (c HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	return keys
}

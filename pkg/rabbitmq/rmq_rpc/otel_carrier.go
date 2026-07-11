package rmqrpc

import amqp "github.com/rabbitmq/amqp091-go"

// TableCarrier adapts amqp.Table to the OpenTelemetry propagation.TextMapCarrier
// interface so trace context can be injected into / extracted from AMQP message headers.
type TableCarrier amqp.Table

// Get returns the value associated with the passed key.
func (c TableCarrier) Get(key string) string {
	v, ok := c[key]
	if !ok {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

// Set stores the key-value pair.
func (c TableCarrier) Set(key, value string) {
	c[key] = value
}

// Keys lists the keys stored in this carrier.
func (c TableCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}

	return keys
}

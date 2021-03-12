package rmq

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/evrone/go-service-template/pkg/logger"
)

type Connection struct {
	url              string
	waitTime         time.Duration
	attempts         int
	consumerExchange string

	connection *amqp.Connection
	channel    *amqp.Channel
	delivery   <-chan amqp.Delivery
}

func newConnection(consumerExchange string) *Connection {
	conn := &Connection{
		url:              "amqp://guest:guest@rabbitmq:5672/",
		waitTime:         5 * time.Second, //nolint:gomnd // will be config
		attempts:         10,              //nolint:gomnd // will be config
		consumerExchange: consumerExchange,
	}

	err := conn.attemptConnect()
	if err != nil {
		logger.Fatal(err, "rmq - newConnection - conn.attemptConnect")
	}

	return conn
}

func (c *Connection) attemptConnect() error {
	var err error
	for i := c.attempts; i > 0; i-- {
		if err = c.connect(); err == nil {
			break
		}

		log.Printf("RabbitMQ is trying to connect, attempts left: %d", i)
		time.Sleep(c.waitTime)
	}

	if err != nil {
		return errors.Wrap(err, "rmq - attemptConnect - c.connect")
	}

	return nil
}

func (c *Connection) connect() error {
	var err error

	c.connection, err = amqp.Dial(c.url)
	if err != nil {
		return errors.Wrap(err, "amqp.Dial")
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		return errors.Wrap(err, "c.connection.Channel")
	}

	err = c.channel.ExchangeDeclare(
		c.consumerExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "c.connection.Channel")
	}

	queue, err := c.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "c.channel.QueueDeclare")
	}

	err = c.channel.QueueBind(
		queue.Name,
		"",
		c.consumerExchange,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "c.channel.QueueBind")
	}

	c.delivery, err = c.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "c.channel.Consume")
	}

	return nil
}

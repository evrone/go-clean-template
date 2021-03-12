package rmq

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Message struct {
	Queue         string
	Priority      uint8
	ContentType   string
	Body          []byte
	ReplyTo       string
	CorrelationID string
}

type pendingCall struct {
	done   chan struct{}
	status string
	body   []byte
}

type Client struct {
	conn           *Connection
	serverExchange string
	timeout        time.Duration
	error          chan error
	stop           chan struct{}

	mx    *sync.RWMutex
	calls map[string]*pendingCall
}

func NewClient(clientExchange, serverExchange string) *Client {
	c := &Client{
		conn:           newConnection(clientExchange),
		serverExchange: serverExchange,
		timeout:        2 * time.Second, //nolint:gomnd // will be config
		error:          make(chan error),
		stop:           make(chan struct{}),
		mx:             &sync.RWMutex{},
		calls:          make(map[string]*pendingCall),
	}

	go c.consumer()

	return c
}

func (c *Client) publish(corrID, handler string, request interface{}) error {
	var (
		requestBody []byte
		err         error
	)

	if request != nil {
		requestBody, err = json.Marshal(request)
		if err != nil {
			return errors.New("json.Marshal")
		}
	}

	err = c.conn.channel.Publish(c.serverExchange, "", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       c.conn.consumerExchange,
			Type:          handler,
			Body:          requestBody,
		})
	if err != nil {
		return errors.Wrap(err, "c.channel.Publish")
	}

	return nil
}

func (c *Client) RemoteCall(handler string, request, response interface{}) error { //nolint:cyclop // complex func
	select {
	case <-c.stop:
		time.Sleep(c.timeout)
		select {
		case <-c.stop:
			return errors.New("rmq - Client - RemoteCall - connection closed")
		default:
		}
	default:
	}

	corrID := uuid.New().String()

	err := c.publish(corrID, handler, request)
	if err != nil {
		return errors.Wrap(err, "rmq - Client - RemoteCall - c.publish")
	}

	call := &pendingCall{done: make(chan struct{})}

	c.addCall(corrID, call)
	defer c.deleteCall(corrID)

	select {
	case <-time.After(c.timeout):
		return ErrTimeout
	case <-call.done:
	}

	if call.status == Success {
		err = json.Unmarshal(call.body, &response)
		if err != nil {
			return errors.Wrap(err, "rmq - Client - RemoteCall - json.Unmarshal")
		}

		return nil
	}

	if call.status == ErrBadHandler.Error() {
		return ErrBadHandler
	}

	if call.status == ErrInternalServer.Error() {
		return ErrInternalServer
	}

	return nil
}

func (c *Client) consumer() {
	for {
		select {
		case <-c.stop:
			return
		case d, opened := <-c.conn.delivery:
			if !opened {
				c.reconnect()

				return
			}

			_ = d.Ack(false) //nolint:errcheck // don't need this

			c.getCall(&d)
		}
	}
}

func (c *Client) reconnect() {
	close(c.stop)

	err := c.conn.attemptConnect()
	if err != nil {
		c.error <- err
		close(c.error)

		return
	}

	c.stop = make(chan struct{})

	go c.consumer()
}

func (c *Client) getCall(d *amqp.Delivery) {
	c.mx.RLock()
	call, ok := c.calls[d.CorrelationId]
	c.mx.RUnlock()

	if !ok {
		return
	}

	call.status = d.Type
	call.body = d.Body
	close(call.done)
}

func (c *Client) addCall(corrID string, call *pendingCall) {
	c.mx.Lock()
	c.calls[corrID] = call
	c.mx.Unlock()
}

func (c *Client) deleteCall(corrID string) {
	c.mx.Lock()
	delete(c.calls, corrID)
	c.mx.Unlock()
}

func (c *Client) Notify() <-chan error {
	return c.error
}

func (c *Client) Shutdown() error {
	select {
	case <-c.error:
		return nil
	default:
	}

	close(c.stop)
	time.Sleep(c.timeout)

	err := c.conn.connection.Close()
	if err != nil {
		return fmt.Errorf("rmq - Client - Shutdown - c.connection.Close: %w", err)
	}

	return nil
}

package gohelprabbitmq

import (
	"errors"
	"time"

	"github.com/streadway/amqp"
)

// AsyncConsumer wraps a SimpleConsumer with a channel to receive the deliveries on demand
type AsyncConsumer struct {
	*SimpleConsumer

	buffer chan *amqp.Delivery
}

// NewAsyncConsumer creates a new BufferedConsumer
func NewAsyncConsumer(connection *Connection, name string) *AsyncConsumer {
	return &AsyncConsumer{
		NewSimpleConsumer(connection, name),
		make(chan *amqp.Delivery),
	}
}

// Consume will start consuming the queue
func (consumer *AsyncConsumer) Consume() (err error) {
	return consumer.SimpleConsumer.Consume(func(delivery amqp.Delivery) {
		consumer.buffer <- &delivery
	})
}

// AwaitDelivery will await for the next delivery or return an error if it did not arrive in time
func (consumer *AsyncConsumer) AwaitDelivery(timeout time.Duration) (delivery *amqp.Delivery, err error) {
	select {
	case <-time.After(timeout):
		return nil, errors.New("Awaiting delivery timed out")

	case delivery = <-consumer.buffer:
		return
	}
}

package gohelprabbitmq

import (
	"strings"

	"github.com/streadway/amqp"
)

// SimplePublisher is a publisher which will simply publish in a queue
type SimplePublisher struct {
	connection *Connection

	exchange string
	key      string

	Mandatory bool //If not set will drop the message if not placed in a queue
	Immediate bool //If set and has no consumer ready the message will be returned
}

// NewSimplePublisher creates a SimplePublisher structure with it's default values
func NewSimplePublisher(connection *Connection, path string) *SimplePublisher {
	splitPath := strings.Split(path, "@")
	exchange := splitPath[0]
	key := ""
	if len(splitPath) > 1 {
		key = splitPath[1]
	}

	return &SimplePublisher{
		connection: connection,
		exchange:   exchange,
		key:        key,

		Mandatory: false,
		Immediate: false,
	}
}

// GetExchange gets the exchange
func (publisher *SimplePublisher) GetExchange() (queueName string) {
	return publisher.exchange
}

// GetKey gets the key
func (publisher *SimplePublisher) GetKey() (queueName string) {
	return publisher.key
}

// Publish publishes a message directly into a queue
func (publisher *SimplePublisher) Publish(message []byte) (err error) {
	channel, err := publisher.connection.Connect()
	if err != nil {
		return
	}
	defer channel.Close()

	channel.Publish(publisher.exchange,
		publisher.key,
		publisher.Mandatory,
		publisher.Immediate,
		amqp.Publishing{
			Body: message,
		},
	)
	return
}

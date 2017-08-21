package gohelprabbitmq

import (
	"github.com/streadway/amqp"
)

// Channel joins a controlled connection with RabbitMQ's channel
type Channel struct {
	*Connection
	*amqp.Channel
}

// Close closes the underlying connection if it's the last connection and the channel itself
func (channel *Channel) Close() (err error) {
	channelError := channel.Channel.Close()
	connectionError := channel.Connection.Close()
	if connectionError != nil {
		return connectionError
	}
	return channelError
}

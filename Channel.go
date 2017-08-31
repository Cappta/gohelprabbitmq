package gohelprabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

// Channel joins a controlled connection with RabbitMQ's channel
type Channel struct {
	*Connection
	*amqp.Channel
}

// ConsumeQueue will consume the queue using the received settings and queue name on the channel itself
func (channel *Channel) ConsumeQueue(settings *ConsumerSettings, queueName string) (deliveryChannel <-chan amqp.Delivery, err error) {
	return channel.Consume(
		queueName,
		settings.name,
		settings.AutoAck,
		settings.Exclusive,
		settings.NoLocal,
		settings.NoWait,
		settings.Arguments,
	)
}

// DeclareQueue declares the queue with the received settings on the channel itself
func (channel *Channel) DeclareQueue(settings *QueueSettings) (queue amqp.Queue, err error) {
	return channel.QueueDeclare(
		settings.name,
		settings.Durable,
		settings.DeleteWhenUnused,
		settings.Exclusive,
		settings.NoWait,
		settings.Arguments,
	)
}

// StopConsumingQueue will stop consuming the queue on the channel itself
func (channel *Channel) StopConsumingQueue(settings *ConsumerSettings) (err error) {
	if settings.name == "" {
		return errors.New("Cannot stop an unnamed consumer")
	}

	return channel.Channel.Cancel(settings.name, settings.NoWait)
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

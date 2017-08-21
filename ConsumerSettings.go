package gohelprabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

// ConsumerSettings contains the data which will be used when creating the consumer
type ConsumerSettings struct {
	name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Arguments amqp.Table
}

// NewConsumerSettings creates a QueueSettings structure with it's default values
func NewConsumerSettings(name string) *ConsumerSettings {
	return &ConsumerSettings{
		name:      name,
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Arguments: nil,
	}
}

// GetName returns the consumer name
func (settings *ConsumerSettings) GetName() (name string) {
	return settings.name
}

// ConsumeQueue will consume the queue using the current settings provided the channel and queue name
func (settings *ConsumerSettings) ConsumeQueue(channel *Channel, queueName string) (deliveryChannel <-chan amqp.Delivery, err error) {
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

// StopConsumingQueue will stop consuming the queue on the specified channel
func (settings *ConsumerSettings) StopConsumingQueue(channel *Channel) (err error) {
	if settings.name == "" {
		return errors.New("Cannot stop an unnamed consumer")
	}

	return channel.Cancel(settings.name, settings.NoWait)
}

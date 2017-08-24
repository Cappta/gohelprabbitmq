package gohelprabbitmq

import (
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

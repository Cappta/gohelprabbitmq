package gohelprabbitmq

import (
	"github.com/streadway/amqp"
)

// QueueSettings contains the data which will be used when creating the queue
type QueueSettings struct {
	name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
	Arguments        amqp.Table
}

// NewQueueSettings creates a QueueSettings structure with it's default values
func NewQueueSettings(name string) *QueueSettings {
	return &QueueSettings{
		name:             name,
		Durable:          false,
		DeleteWhenUnused: true,
		Exclusive:        false,
		NoWait:           false,
		Arguments:        nil,
	}
}

// GetName returns the queue name
func (settings *QueueSettings) GetName() (name string) {
	return settings.name
}

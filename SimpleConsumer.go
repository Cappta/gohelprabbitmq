package gohelprabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

// SimpleConsumer is a consumer which will simply consume a queue
type SimpleConsumer struct {
	Connection       *Connection
	QueueSettings    *QueueSettings
	ConsumerSettings *ConsumerSettings

	consumingChannel *Channel
}

// NewSimpleConsumer creates a SimpleConsumer structure
func NewSimpleConsumer(connection *Connection, name string) *SimpleConsumer {
	return &SimpleConsumer{
		connection,
		NewQueueSettings(name),
		NewConsumerSettings(name),
		nil,
	}
}

// Consume connects to RabbitMQ and consumes Queue with the provided callback
func (consumer *SimpleConsumer) Consume(callback func(delivery amqp.Delivery)) (err error) {
	consumer.consumingChannel, err = consumer.Connection.Connect()
	if err != nil {
		return
	}
	defer consumer.closeChannel()

	queue, err := consumer.consumingChannel.DeclareQueue(consumer.QueueSettings)
	if err != nil {
		return
	}

	deliveryChannel, err := consumer.consumingChannel.ConsumeQueue(consumer.ConsumerSettings, queue.Name)
	if err != nil {
		return
	}

	for delivery := range deliveryChannel {
		go callback(delivery)
	}
	return errors.New("Disconnected")
}

func (consumer *SimpleConsumer) closeChannel() {
	consumer.consumingChannel.Close()
	consumer.consumingChannel = nil
}

// StopConsuming will stop the current Consumer
func (consumer *SimpleConsumer) StopConsuming() (err error) {
	if consumer.consumingChannel == nil {
		return errors.New("Not currently consuming")
	}

	return consumer.consumingChannel.StopConsumingQueue(consumer.ConsumerSettings)
}

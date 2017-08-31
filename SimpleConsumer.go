package gohelprabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

// SimpleConsumer is a consumer which will simply consume a queue
type SimpleConsumer struct {
	connection *Connection

	queueSettings *QueueSettings

	consumerSettings *ConsumerSettings

	consumingChannel *Channel
}

// NewSimpleConsumer creates a SimpleConsumer structure
func NewSimpleConsumer(connection *Connection, queueSettings *QueueSettings, consumerSettings *ConsumerSettings) *SimpleConsumer {
	return &SimpleConsumer{
		connection:       connection,
		queueSettings:    queueSettings,
		consumerSettings: consumerSettings,
	}
}

// GetConnection returns the connection struct
func (consumer *SimpleConsumer) GetConnection() (connection *Connection) {
	return consumer.connection
}

// GetQueueSettings returns the queueSettings struct
func (consumer *SimpleConsumer) GetQueueSettings() (queueSettings *QueueSettings) {
	return consumer.queueSettings
}

// GetConsumerSettings returns the consumerSettings struct
func (consumer *SimpleConsumer) GetConsumerSettings() (consumerSettings *ConsumerSettings) {
	return consumer.consumerSettings
}

// Consume connects to RabbitMQ and consumes Queue with the provided callback
func (consumer *SimpleConsumer) Consume(handle func(delivery amqp.Delivery)) (err error) {
	consumer.consumingChannel, err = consumer.connection.Connect()
	if err != nil {
		return
	}
	defer consumer.closeChannel()

	queue, err := consumer.consumingChannel.DeclareQueue(consumer.queueSettings)
	if err != nil {
		return
	}

	deliveryChannel, err := consumer.consumingChannel.ConsumeQueue(consumer.consumerSettings, queue.Name)
	if err != nil {
		return
	}

	for delivery := range deliveryChannel {
		go handle(delivery)
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

	return consumer.consumingChannel.StopConsumingQueue(consumer.consumerSettings)
}

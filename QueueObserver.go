package gohelprabbitmq

import (
	"errors"
	"time"
)

// QueueObserver observes a queue on RabbitMQ's server
type QueueObserver struct {
	connection *Connection
	settings   *QueueSettings
}

// NewQueueObserver creates a new QueueObserver
func NewQueueObserver(connection *Connection, settings *QueueSettings) *QueueObserver {
	return &QueueObserver{
		connection: connection,
		settings:   settings,
	}
}

// AwaitConsumer will await an consumer for the observed queue
func (observer *QueueObserver) AwaitConsumer(timeout time.Duration) (err error) {
	channel, err := observer.connection.Connect()
	if err != nil {
		return
	}
	defer channel.Close()

	timeoutUnixNano := time.Now().Add(timeout).UnixNano()
	for time.Now().UnixNano() <= timeoutUnixNano {
		queue, err := channel.DeclareQueue(observer.settings)
		if err != nil {
			return err
		}

		if queue.Consumers > 0 {
			return nil
		}
	}
	return errors.New("Awaiting consumer timed out")
}

package gohelprabbitmq

import (
	"errors"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// TemporaryConsumer is a consumer that does not last forever
type TemporaryConsumer struct {
	*SimpleConsumer
}

// NewTemporaryConsumer creates a new TemporaryConsumer
func NewTemporaryConsumer(simpleConsumer *SimpleConsumer) *TemporaryConsumer {
	return &TemporaryConsumer{
		simpleConsumer,
	}
}

// AwaitDelivery will await the delivery for the specified timeout
func (consumer *TemporaryConsumer) AwaitDelivery(timeout time.Duration) (delivery *amqp.Delivery, err error) {
	done := make(chan bool)
	once := &sync.Once{}
	stopTimer := time.AfterFunc(timeout, func() {
		once.Do(func() {
			err = consumer.SimpleConsumer.StopConsuming()
			if err == nil {
				err = errors.New("Awaiting delivery timed out")
			}
			done <- true
		})
	})
	consumer.SimpleConsumer.Consume(func(consumedDelivery amqp.Delivery) {
		once.Do(func() {
			stopTimer.Stop()
			delivery = &consumedDelivery
			err = consumer.SimpleConsumer.StopConsuming()
			done <- true
		})
	})
	<-done
	close(done)
	return
}

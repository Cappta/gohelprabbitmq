package gohelprabbitmq

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

var (
	errContainerMissingRPC = errors.New("JSON lacks RPC data Content")
	errCallbackExpired     = errors.New("Callback expired")
)

const (
	rpcPath      = "RPC"
	rpcIDPath    = "RPC.ID"
	rpcQueuePath = "RPC.Queue"
)

// RPC is a remote procedure call concentrator that uses a single rabbitmq queue a single json message
type RPC struct {
	*SimpleConsumer
	mutex       *sync.Mutex
	callbackMap map[string]chan *amqp.Delivery
	timeout     time.Duration
}

// NewRPC creates a new RPC
func NewRPC(consumer *SimpleConsumer, timeout time.Duration) *RPC {
	return &RPC{
		consumer,
		&sync.Mutex{},
		make(map[string]chan *amqp.Delivery),
		timeout,
	}
}

// Prepare prepares a container given a consumer to be able to forward the message and receive the reply
func (rpc *RPC) Prepare(container *gabs.Container) <-chan *amqp.Delivery {
	id := uuid.NewV4().String()
	rpc.appendCallbackData(id, container)
	return rpc.allocateCallback(id)
}

func (rpc *RPC) appendCallbackData(id string, container *gabs.Container) {
	containerParent := getRPCContainerParent(container)
	containerParent.SetP(id, rpcIDPath)
	containerParent.SetP(fmt.Sprintf("@%s", rpc.QueueSettings.GetName()), rpcQueuePath)
}

func (rpc *RPC) allocateCallback(id string) <-chan *amqp.Delivery {
	rpc.mutex.Lock()
	defer rpc.mutex.Unlock()

	channel := make(chan *amqp.Delivery)
	rpc.callbackMap[id] = channel
	time.AfterFunc(rpc.timeout, func() {
		rpc.mutex.Lock()
		defer rpc.mutex.Unlock()

		close(channel)
		delete(rpc.callbackMap, id)
	})
	return channel
}

// Consume consumes the SimpleConsumer as a source to the RPC response
func (rpc *RPC) Consume() error {
	return rpc.SimpleConsumer.Consume(func(delivery amqp.Delivery) {
		err := rpc.routeDelivery(&delivery)
		if err != nil {
			log.Println("Error parsing JSON \"", err, "\" Content: ", string(delivery.Body))
			delivery.Reject(false)
		}
	})
}

func (rpc *RPC) routeDelivery(delivery *amqp.Delivery) (err error) {
	container, err := gabs.ParseJSON(delivery.Body)
	if err != nil {
		return
	}

	callback, err := rpc.popDeliveryCallback(container)
	if err != nil {
		return
	}

	delivery.Body = container.Bytes()
	callback <- delivery
	return
}

func (rpc *RPC) popDeliveryCallback(container *gabs.Container) (callback chan *amqp.Delivery, err error) {
	rpcContainerParent := getRPCContainerParent(container)
	if rpcContainerParent.ExistsP(rpcIDPath) == false {
		return nil, errContainerMissingRPC
	}
	id := rpcContainerParent.Path(rpcIDPath).Data().(string)
	err = rpcContainerParent.Delete(rpcPath)
	if err != nil {
		return
	}

	rpc.mutex.Lock()
	defer rpc.mutex.Unlock()
	if callback, ok := rpc.callbackMap[id]; ok {
		return callback, nil
	}
	return nil, errCallbackExpired
}

func popRPCQueue(container *gabs.Container) (queue string, err error) {
	rpcContainerParent := getRPCContainerParent(container)
	if rpcContainerParent.ExistsP(rpcIDPath) == false {
		return "", errContainerMissingRPC
	}
	queue = rpcContainerParent.Path(rpcQueuePath).Data().(string)
	rpcContainerParent.DeleteP(rpcQueuePath)
	return
}

func getRPCContainerParent(container *gabs.Container) (rpcContainer *gabs.Container) {
	rpcContainer = container
	for rpcContainer.Exists(rpcPath, rpcPath) {
		rpcContainer = rpcContainer.Search(rpcPath)
	}
	return
}

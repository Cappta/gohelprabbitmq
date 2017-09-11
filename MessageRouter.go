package gohelprabbitmq

import (
	"github.com/Jeffail/gabs"
)

// MessageRouter routes the message through a rabbitmq connection given a primary forward to path or return an existing RPC call
type MessageRouter struct {
	forwardToPath string
	connection    *Connection
}

// NewMessageRouter creates a new MessageRouter
func NewMessageRouter(forwardToPath string, connection *Connection) (messageRouter *MessageRouter) {
	return &MessageRouter{
		forwardToPath,
		connection,
	}
}

// Route will figute out the container's next path and route it accordingly
func (messageRouter *MessageRouter) Route(container *gabs.Container) (err error) {
	path, err := messageRouter.popRoutePath(container)
	if err != nil {
		return
	}
	publisher := NewSimplePublisher(messageRouter.connection, path)
	return publisher.Publish(container.Bytes())
}

func (messageRouter *MessageRouter) popRoutePath(container *gabs.Container) (routePath string, err error) {
	if container.ExistsP(messageRouter.forwardToPath) == false {
		return popRPCQueue(container)
	}

	routePath = container.Path(messageRouter.forwardToPath).Data().(string)
	container.DeleteP(messageRouter.forwardToPath)
	return
}

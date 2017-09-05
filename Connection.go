package gohelprabbitmq

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

// Connection contains the Connection data
type Connection struct {
	user     string
	password string
	host     string
	port     int

	connection *amqp.Connection

	activeConnections int
	mutex             *sync.Mutex
}

// NewConnection creates a Connection structure with it's default values
func NewConnection(user, password, host string, port int) *Connection {
	return &Connection{
		user:     user,
		password: password,

		host: host,
		port: port,

		connection: nil,

		activeConnections: 0,
		mutex:             new(sync.Mutex),
	}
}

// ConnectLocallyWithDefaultUser will return a new Connection with the guest user and password on localhost at the default port 5672
func ConnectLocallyWithDefaultUser() *Connection {
	return NewConnection("guest", "guest", "127.0.0.1", 5672)
}

// GetConnection returns the current Connection
func (connection *Connection) GetConnection() *amqp.Connection {
	return connection.connection
}

// Connect will connect if disconnected, increment the active connections and return a new channel
func (connection *Connection) Connect() (channel *Channel, err error) {
	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	connection.activeConnections++
	if connection.activeConnections == 1 {
		connection.connection, err = amqp.Dial(connection.getConnectionString())

		if err != nil {
			connection.Close()
			return
		}
	}

	amqpChannel, err := connection.connection.Channel()
	if err != nil {
		connection.Close()
	}
	channel = &Channel{connection, amqpChannel}
	return
}

func (connection *Connection) getConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", connection.user, connection.password, connection.host, connection.port)
}

// Close will decrease the Connection reference and if nobody is using the Connection will disconnect from the server
func (connection *Connection) Close() (err error) {
	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	connection.activeConnections--
	if connection.activeConnections == 0 {
		err = connection.connection.Close()
	}
	return
}

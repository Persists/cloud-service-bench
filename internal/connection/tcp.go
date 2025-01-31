package connection

import (
	"fmt"
	"net"
)

type ConnectionClient struct {
	connection *net.Conn
	address    string
}

func NewConnectionClient(host string, port int) *ConnectionClient {
	return &ConnectionClient{
		address: fmt.Sprintf("%s:%d", host, port),
	}
}

// Connect establishes a TCP connection to the server.
func (cc *ConnectionClient) Connect() error {
	conn, err := net.Dial("tcp", cc.address)
	if err != nil {
		return err
	}

	cc.connection = &conn
	return nil
}

func (cc *ConnectionClient) Disconnect() error {
	err := (*cc.connection).Close()
	if err != nil {
		return err
	}

	return nil
}

// SendMessage sends a message over the TCP connection.
// It takes a string message as input and returns an error if the message
// could not be sent.
//
// Parameters:
//   - message: The string message to be sent.
func (cc *ConnectionClient) SendMessage(message string) error {
	if _, err := (*cc.connection).Write([]byte(message)); err != nil {
		return err
	}

	return nil
}

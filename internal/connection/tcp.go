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

func (cc *ConnectionClient) SendMessage(message string) error {
	if _, err := (*cc.connection).Write([]byte(message)); err != nil {
		return err
	}

	return nil
}

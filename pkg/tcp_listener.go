package blaster

import (
	"fmt"
	"net"
)

type ListenerOptions func(t *TcpListener)

func WithPort(port int) ListenerOptions {
	return func(l *TcpListener) {
		l.Port = port
	}
}

type TcpListener struct {
	Port int
}

func NewTcpListener(opts ...ListenerOptions) *TcpListener {
	l := &TcpListener{
		Port: 9293,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *TcpListener) Listen() error {
	PORT := fmt.Sprintf(":%d", l.Port)
	listener, err := net.Listen("tcp4", PORT)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}

		go l.handleConnection(c)
	}
}

func (l *TcpListener) handleConnection(c net.Conn) {

}

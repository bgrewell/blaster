package schedulers

import (
	blaster "github.com/BGrewell/blaster/internal"
	"net"
)

type Scheduler interface {
	Identifier() string
	Handle(c net.Conn, params *blaster.TcpFlow, cancel <-chan interface{})
}

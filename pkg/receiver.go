package blaster

import (
	"bufio"
	"github.com/BGrewell/blaster/internal"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"sync"
)

func HandleRecv(c net.Conn, flow *internal.TcpFlow, wg *sync.WaitGroup) {
	defer wg.Done()

	// TODO: Plumb so that this can actually bue used to cancel a flow
	cancel := make(chan interface{})
	recvBuff := make([]byte, flow.PacketSize*2)
	reader := bufio.NewReader(c)

	for {
		select {
		case <- cancel:
			return
		default:
			read, err := reader.Read(recvBuff)
			if err == io.EOF {
				return
			} else if err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"client": c.RemoteAddr().String(),
					"read": read,
				}).Trace("failed to read from tcp connection")
				return
			}
			// TODO: Update accounting
		}
	}
}

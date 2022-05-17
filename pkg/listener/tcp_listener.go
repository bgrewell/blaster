package listener

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/BGrewell/blaster/pkg"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

const (
	DELIM = 0xFF
)

type TcpListener struct {
	Port int `json:"port,omitempty"`
}

func NewTcpListener(opts ...Options) *TcpListener {
	l := &TcpListener{
		Port: 9293,
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l *TcpListener) Listen(readyIndicator chan interface{}) error {
	PORT := fmt.Sprintf(":%d", l.Port)
	listener, err := net.Listen("tcp4", PORT)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.WithFields(log.Fields{
		"port": PORT,
	}).Trace("listening for new connections")
	if readyIndicator != nil {
		readyIndicator <- true
	}

	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"client": c.RemoteAddr().String(),
		}).Trace("accepted new tcp connection")
		go l.handleConnection(c)
	}
}

func (l *TcpListener) handleConnection(c net.Conn) {
	log.WithFields(log.Fields{
		"client": c.RemoteAddr().String(),
	}).Trace("servicing new client")

	for {
		// Do connection setup
		data, err := bufio.NewReader(c).ReadBytes(DELIM)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"client": c.RemoteAddr().String(),
			}).Error("failed to read from connection")
			return
		}

		var params blaster.SessionParameters
		err = json.Unmarshal(data[:len(data)-1], &params)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"client": c.RemoteAddr().String(),
			}).Error("failed to unmarshal session parameters")
		}

		active := sync.WaitGroup{}
		if params.UplinkFlow != nil {
			// Setup receiver for uplink flow
			go blaster.HandleRecv(c, params.UplinkFlow, &active)
			active.Add(1)
			log.WithFields(log.Fields{
				"client": c.RemoteAddr().String(),
			}).Trace("started new data receiver to accept uplink traffic")
		}

		if params.DownlinkFlow != nil {
			// Setup sender for downlink flow
			go blaster.HandleSend(c, params.DownlinkFlow, &active)
			active.Add(1)
			log.WithFields(log.Fields{
				"client": c.RemoteAddr().String(),
			}).Trace("started new data sender to transmit downlink traffic")
		}

		// wait for flow to finish then return
		active.Wait()
		log.WithFields(log.Fields{
			"client": c.RemoteAddr().String(),
		}).Trace("transmissions have finished")
		return
	}
}

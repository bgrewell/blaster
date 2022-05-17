package client

import (
	"encoding/json"
	"fmt"
	"github.com/BGrewell/blaster/internal"
	blaster "github.com/BGrewell/blaster/pkg"
	"github.com/BGrewell/blaster/pkg/listener"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

type TcpClient struct {
	ServerAddr string
	DownlinkFlow *internal.TcpFlow
	UplinkFlow *internal.TcpFlow
}

func NewTcpClient(opts ...Options) *TcpClient {
	c := &TcpClient{
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (tc *TcpClient) RunWithHost(host string, port int, session *blaster.SessionParameters) error {
	og := tc.ServerAddr
	tc.ServerAddr = fmt.Sprintf("%s:%d", host, port)
	err := tc.Run(session)
	tc.ServerAddr = og
	return err
}

func (tc *TcpClient) Run(session *blaster.SessionParameters) error {
	// Connect to listener
	c, err := net.Dial("tcp", tc.ServerAddr)
	if err != nil {
		return err
	}

	// Send parameters
	params, err := json.Marshal(session)
	if err != nil {
		return err
	}
	buff := make([]byte, len(params) + 1)
	copy(buff, params)
	buff[len(buff) - 1] = listener.DELIM
	_, err = c.Write(buff)
	if err != nil {
		return err
	}

	active := sync.WaitGroup{}
	if session.DownlinkFlow != nil {
		go blaster.HandleRecv(c, session.DownlinkFlow, &active)
		active.Add(1)
		log.WithFields(log.Fields{
			"server": c.RemoteAddr().String(),
		}).Trace("started new data receiver to accept uplink traffic")
	}

	if session.UplinkFlow != nil {
		go blaster.HandleSend(c, session.UplinkFlow, &active)
		active.Add(1)
		log.WithFields(log.Fields{
			"server": c.RemoteAddr().String(),
		}).Trace("started new data sender to transmit downlink traffic")
	}

	// wait for flow to finish then return
	active.Wait()
	log.WithFields(log.Fields{
		"server": c.RemoteAddr().String(),
	}).Trace("transmissions have finished")

	return nil
}
package client

import (
	"fmt"
	"github.com/BGrewell/blaster/internal"
)

type Options func(t *TcpClient)

func WithServerAddr(host string, port int) Options {
	return func(c *TcpClient) {
		c.ServerAddr = fmt.Sprintf("%s:%d", host, port)
	}
}

func WithUplinkFlow(flow *internal.TcpFlow) Options {
	return func(c *TcpClient) {
		c.UplinkFlow = flow
	}
}

func WithDownlinkFlow(flow *internal.TcpFlow) Options {
	return func(c *TcpClient) {
		c.DownlinkFlow = flow
	}
}
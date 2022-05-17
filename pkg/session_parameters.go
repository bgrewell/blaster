package blaster

import "github.com/BGrewell/blaster/internal"

type SessionParameters struct {
	UplinkFlow   *internal.TcpFlow `json:"uplink_flow"`
	DownlinkFlow *internal.TcpFlow `json:"downlink_flow"`
}

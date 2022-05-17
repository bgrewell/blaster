package schedulers

import (
	blaster "github.com/BGrewell/blaster/internal"
	log "github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net"
	"time"
)

type TimedScheduler struct {
}

func (ts *TimedScheduler) Identifier() string {
	return "timed"
}

func (ts *TimedScheduler) Handle(c net.Conn, flow *blaster.TcpFlow, cancel <- chan interface{}) {

	payload := make([]byte, flow.PacketSize * 100)
	rand.Seed(time.Now().UnixNano())
	rand.Read(payload)
	pmax := flow.PacketSize * 99

	payloadSize := flow.PacketSize
	packetsPerSecond := float64(flow.RateBitsPerSec) / float64(payloadSize)
	interval := 1 / packetsPerSecond
	log.WithFields(log.Fields{
		"payload": payloadSize,
		"rate": flow.RateBitsPerSec,
		"interval": interval,
	}).Trace("calculated timing")
	// Wait till start
	log.Trace("waiting for start time")
	for flow.StartTime > time.Now().UnixNano() {
		time.Sleep(1 * time.Microsecond)
	}

	stop := make(chan interface{})
	// Setup stop channel
	go func() {
		<- time.After(time.Duration(flow.Duration))
		stop <- true
	}()
	log.Trace("setting up test stop")

	log.Trace("starting test")
	next := time.Now()
	for {
		select {
		case <- stop:
			// Stop time hit. Return to stop sending
			c.Close()
			return
		case <- cancel:
			// Sending was canceled
			return
		default:
			// Send at timed rate
			idx := rand.Intn(pmax)
			for time.Now().Before(next){
				//TODO: A better approach is to use sleeps = 1/4 of the time remaining if it is over 100 micro seconds then spin
				//TODO: this gets just as good precision but uses much less CPU resources
			}
			sent, err := c.Write(payload[idx : idx+flow.PacketSize])
			if err == io.EOF {
				log.WithFields(log.Fields{
					"err": err,
					"client": c.RemoteAddr().String(),
				}).Debug("connection closed")
				return
			} else if err != nil {
				log.WithFields(log.Fields{
					"err": err,
					"sent": sent,
				}).Debug("failed to send payload")
				return
			}
			next = time.Now().Add(time.Duration(interval) * time.Second)
			// TODO: Update accounting
		}
	}
}

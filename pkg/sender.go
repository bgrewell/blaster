package blaster

import (
	"github.com/BGrewell/blaster/internal"
	"github.com/BGrewell/blaster/internal/schedulers"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

func HandleSend(c net.Conn, flow *internal.TcpFlow, wg *sync.WaitGroup) {
	defer wg.Done()
	var scheduler schedulers.Scheduler

	switch flow.Scheduler {
	case "token":
		scheduler = &schedulers.TokenScheduler{}
	case "timed":
		scheduler = &schedulers.TimedScheduler{}
	case "soak":
		scheduler = &schedulers.SoakScheduler{}
	default:
		log.WithFields(log.Fields{
			"client": c.RemoteAddr().String(),
			"scheduler": flow.Scheduler,
			"flow": flow,
		}).Error("unknown scheduler")
		return
	}

	// TODO: Plumb so that this can actually be used to cancel a flow
	cancel := make(chan interface{})
	scheduler.Handle(c, flow, cancel)

	return
}

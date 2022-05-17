package main

import (
  "fmt"
  "github.com/BGrewell/blaster/internal"
  blaster "github.com/BGrewell/blaster/pkg"
  "github.com/BGrewell/blaster/pkg/client"
  "github.com/BGrewell/blaster/pkg/listener"
  log "github.com/sirupsen/logrus"
  "sync"
  "time"
)

var (
  l *listener.TcpListener
)

func RunBiDi(wg *sync.WaitGroup) {
  defer wg.Done()
  ready := make(chan interface{})

  if l == nil {
    l = listener.NewTcpListener(listener.WithPort(9006))
    go func() {
      // Need a good way to cancel the listener
      err := l.Listen(ready)
      if err != nil {
        log.WithFields(log.Fields{
          "err": err,
        }).Fatal("failed to start tcp listener")
      }
    }()
    <- ready
  }

  flow := &internal.TcpFlow{
    StartTime:      time.Now().UnixNano() + (5 * time.Second).Nanoseconds(),
    Duration:       (30 * time.Second).Nanoseconds(),
    PacketSize:     1000,
    RateBitsPerSec: 50 * 1000 * 1000,
    Scheduler:      "soak",
  }

  session := &blaster.SessionParameters{
    UplinkFlow:   flow,
    DownlinkFlow: flow,
  }

  c := client.NewTcpClient(client.WithServerAddr("127.0.0.1", 9006))
  err := c.Run(session)
  if err != nil {
    log.WithFields(log.Fields{
      "err": err,
    }).Fatal("failed to run tcp client")
  }
}

func main() {

  log.SetLevel(log.TraceLevel)

  fmt.Println("[+] Starting Listener on port 9006")

  wg := sync.WaitGroup{}
  wg.Add(2)

  go RunBiDi(&wg)
  time.Sleep(10 * time.Second)
  go RunBiDi(&wg)

  wg.Wait()
}

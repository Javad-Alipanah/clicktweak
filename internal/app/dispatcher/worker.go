package dispatcher

import (
	"clicktweak/internal/pkg/model"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

// Workers represents access logging worker
type Workers struct {
	// sockAddr is the address of log forwarder
	sockAddr string

	// log is the channel in which logs are received
	log chan model.Log
}

func NewWorkers(addr string, log chan model.Log) *Workers {
	return &Workers{addr, log}
}

func (w *Workers) connectAndHandle(i int) {
	for {
		c, err := net.Dial("tcp", w.sockAddr)
		if err != nil {
			log.Error(err)
			time.Sleep(time.Millisecond * 10)
			continue
		}
		if !w.collectAndSend(c, i) {
			break
		}
	}
}

// collects access log and forwards the log
func (w *Workers) collectAndSend(c net.Conn, i int) bool {
	defer c.Close()

LOOP:
	for {
		select {
		// send json log to fluent-bit
		case logElem, more := <-w.log:
			// workers closed
			if !more {
				return false
			}

			logJSON, err := json.Marshal(&logElem)
			if err != nil {
				log.Errorf("worker %d: %s\n", i, err)
				continue LOOP
			}

			logJSON = append(logJSON, '\n')
			n, err := c.Write(logJSON)
			if err != nil {
				// connection closed
				if err == io.EOF {
					log.Error("log forwarder connection closed unexpectedly")
					return true
					// other errors
				} else {
					log.Error("log forwarder write error: ", err.Error())
					return true
				}
			}

			// data transmission failed partially
			if n != len(logJSON) {
				log.Error("write to log forwarder returned ", n, " instead of ", len(logJSON))
			}
		}
	}
}

// Run starts workers in separate goroutines to collectAndHandle logs
//
// `n` is the number of workers
func (w *Workers) Run(n int) {
	for i := 0; i < n; i++ {
		go w.connectAndHandle(i)
	}
}

// Close stops all workers
func (w *Workers) Close() {
	close(w.log)
}

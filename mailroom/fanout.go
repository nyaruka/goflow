package mailroom

import (
	"fmt"
	"time"

	"github.com/nyaruka/goflow/mailroom/store"
)

// Starts our goroutine popping msgs off our redis queue and pushing them onto our
// handler queue. We only pop works off when we have at least one available worker
func startFanout(m *mailroom) {
	// create our workers, put them on our worker channel
	for i := 0; i < m.config.Workers; i++ {
		newWorker(i+1, m).start()
	}

	go func() {
		m.waitGroup.Add(1)

		for {
			select {
			case worker := <-m.workerChan:
				// we got an available worker, get our next message
				msg, err := m.popNextMsg()
				if err != nil && msg != nil {
					worker.msgChan <- msg
				} else {
					// TODO: use subscribe on our list to make things smarter
					time.Sleep(time.Second)
					fmt.Printf("No work, sleeping\n")
				}

			case <-m.stopChan:
				m.waitGroup.Done()
				return
			}
		}
	}()
}

type worker struct {
	id       int
	mailroom *mailroom
	msgChan  chan *store.Msg
}

func newWorker(id int, m *mailroom) *worker {
	msgChan := make(chan *store.Msg)
	return &worker{
		id:       id,
		mailroom: m,
		msgChan:  msgChan,
	}
}

func (w *worker) start() {
	go func() {
		w.mailroom.waitGroup.Add(1)

		// mark ourselves as ready
		w.mailroom.workerChan <- w

		fmt.Printf("Started worker: %d\n", w.id)

		for {
			// and wait for a job
			select {
			case msg := <-w.msgChan:
				fmt.Printf("handling: %v\n", msg)
				time.Sleep(10 * time.Second)
				w.mailroom.workerChan <- w

			case <-w.mailroom.stopChan:
				w.mailroom.waitGroup.Done()
				return
			}
		}
	}()
}

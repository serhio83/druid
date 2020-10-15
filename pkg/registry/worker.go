package registry

import (
	"fmt"
	"log"
	"time"

	"github.com/serhio83/druid/pkg/config"
	u "github.com/serhio83/druid/pkg/utils"
	"go.etcd.io/bbolt"
)

const (
	logHeader = "[druid_worker]"
)

// Worker ...
type Worker struct {
	*config.Config
	*time.Ticker
	repos *Repos
	*bbolt.DB
	stop chan bool
}

// NewWorker creates new registry worker
func NewWorker(interval time.Duration, c *config.Config, db *bbolt.DB, stop <-chan bool) *Worker {
	w := new(Worker)
	w.DB = db
	w.Config = c
	w.Ticker = time.NewTicker(interval)
	log.Println(u.Envelope(logHeader + " docker registry worker started"))
	go w.process(stop)
	return w
}

// process registry images
func (w *Worker) process(stop <-chan bool) {
	for {
		select {
		case <-w.C:
			images, err := ListImages(w.Config)
			if err != nil {
				log.Fatal(err)
			}

			w.repos = images
			for _, rep := range w.repos.Repositories[:] {
				if err := w.storeBucket(string(rep)); err != nil {
					log.Fatal(u.Envelope(fmt.Sprintf("cant store bucket: %v", err)))
				}
			}
			log.Println(u.Envelope(fmt.Sprintf("%s images processed: %d", logHeader, len(w.repos.Repositories))))
		case <-stop:
			log.Println(u.Envelope(logHeader + " worker stopped"))
			return
		}
	}
}

func (w *Worker) storeBucket(name string) error {
	// Start a write transaction.
	if err := w.DB.Update(func(tx *bbolt.Tx) error {
		// Create a bucket.
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			log.Println(u.Envelope(fmt.Sprintf("fail in tx: %v", err)))
			return err
		}
		return nil
	}); err != nil {
		log.Println(u.Envelope(fmt.Sprintf("cant create bucket: %v", err)))
		return err
	}
	return nil
}

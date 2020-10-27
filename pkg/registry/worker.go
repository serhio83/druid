package registry

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
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
	stop    chan bool
	tagProc *TagProcessor
}

// TagProcessor struct
type TagProcessor struct {
	count    int64
	maxCount int64
	*time.Ticker
	tags     []Tag
	image    chan string
	syncChan chan struct{}
	syncFlag bool
	stop     chan bool
	db       *bbolt.DB
	cfg      *config.Config
}

// NewWorker creates new registry worker
func NewWorker(interval time.Duration, c *config.Config, db *bbolt.DB, stop <-chan bool, wg *sync.WaitGroup) *Worker {

	w := new(Worker)
	w.DB = db
	w.Config = c
	w.Ticker = time.NewTicker(interval)
	log.Println(u.Envelope(logHeader + " docker registry worker started"))

	stopProc := make(chan bool, 1)
	imgChan := make(chan string, 10000)

	go func() {
		w.tagProc = newTagProcessor(3000, stopProc, imgChan, w.DB, c)
	}()

	go w.process(stop, stopProc, wg)

	return w
}

// process registry images
func (w *Worker) process(stop <-chan bool, stopProc chan<- bool, wg *sync.WaitGroup) {
	for {
		select {
		case <-w.C:
			log.Println(u.Envelope(fmt.Sprintf("%s [w.tick] start processing on new tick", logHeader)))
			images, err := ListImages(w.Config)
			if err != nil {
				log.Fatal(err)
			}

			w.repos = images

			// we can regulate how many images to process in this slice w.repos.Repositories[:5]
			for _, rep := range w.repos.Repositories[:] {
				if err := w.storeBucket(string(rep)); err != nil {
					log.Fatal(u.Envelope(fmt.Sprintf("cant store bucket: %v", err)))
				}

				// send image name to tags processor`s channel
				w.tagProc.send(string(rep))
			}

			log.Println(u.Envelope(fmt.Sprintf("%s [w.send] total images: %d", logHeader, len(w.repos.Repositories[:]))))

		case <-stop:
			defer wg.Done()
			stopProc <- true
			log.Println(u.Envelope(logHeader + " [w.process] is stopped"))
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

func (tp *TagProcessor) imageSave(image, tag, date string) error {

	if err := tp.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(image))
		err := b.Put([]byte(tag), []byte(date))
		if err != nil {
			log.Println(u.Envelope(fmt.Sprintf("%v", err)))
			return err
		}
		return nil

	}); err != nil {
		log.Println(u.Envelope(fmt.Sprintf("cant update bucket: %v", err)))
		return err
	}
	return nil
}

func newTagProcessor(maxCount int64, stop <-chan bool, image chan string, db *bbolt.DB, c *config.Config) *TagProcessor {

	tp := new(TagProcessor)
	tp.Ticker = time.NewTicker(time.Minute)
	tp.maxCount = maxCount
	tp.image = image
	tp.syncChan = make(chan struct{}, 1)
	tp.syncFlag = false
	tp.db = db
	tp.cfg = c
	go tp.clean()
	go tp.process()

	return tp
}

func (tp *TagProcessor) send(img string) {
	tp.image <- img
	// log.Println(u.Envelope(fmt.Sprintf("%s [t.send] image: %s", logHeader, img)))
}

func (tp *TagProcessor) clean() {
	for {
		<-tp.C
		atomic.StoreInt64(&tp.count, 0)
		if tp.syncFlag {
			tp.syncChan <- struct{}{}
			tp.syncFlag = false
		}
	}
}

// process with rate limit
func (tp *TagProcessor) process() {
	for {
		img := <-tp.image
		if atomic.LoadInt64(&tp.count) >= tp.maxCount {
			tp.syncFlag = true
			<-tp.syncChan
		}

		tl := ListTags(tp.cfg, img)
		// log.Println(u.Envelope(fmt.Sprintf("%s [t.receive] image: %s tags count: %d", logHeader, tl.Name, len(tl.Tags))))

		for _, tg := range tl.Tags {
			go func(tg Tag) {
				t := string(tg)
				d := GetCreationDate(tp.cfg, img, t)
				tp.imageSave(img, t, d)

				// log.Println(u.Envelope(fmt.Sprintf("%s [t.receive] image: %s:%s %v", logHeader, img, t, d)))
			}(tg)
		}
		atomic.AddInt64(&tp.count, 1)
	}
}

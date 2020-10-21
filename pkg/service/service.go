package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	c "github.com/serhio83/druid/pkg/config"
	h "github.com/serhio83/druid/pkg/handlers"
	w "github.com/serhio83/druid/pkg/registry"
	u "github.com/serhio83/druid/pkg/utils"
	v "github.com/serhio83/druid/pkg/version"
	"go.etcd.io/bbolt"
)

const (
	logHeader = "[druid_main]"
)

// Service struct
type Service struct {
	Config *c.Config
	DB     *bbolt.DB
	*w.Worker
}

// Initialize service
func (s *Service) Initialize() {

	s.Config = c.New()
	db, err := bbolt.Open(s.Config.BboltPath+"/bbolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	s.DB = db

	log.Println(u.Envelope(fmt.Sprintf(
		"%s Configured registry: %s:%s",
		logHeader, s.Config.RegHost, s.Config.RegPort,
	)))

	log.Println(u.Envelope(fmt.Sprintf(
		"%s Bbolt storage: %s",
		logHeader, s.Config.BboltPath+"/bbolt.db",
	)))
}

// Run service
func (s *Service) Run() {
	defer s.DB.Close()

	log.Println(u.Envelope(fmt.Sprintf(
		"%s Starting the druid. Commit: %s, build time: %s, release: %s",
		logHeader, v.Commit, v.BuildTime, v.Release,
	)))

	// setup signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// api server port setup
	port := s.Config.ListenPort
	if port == "" {
		log.Fatal(u.Envelope(fmt.Sprintf("%s Port is not set", logHeader)))
	}

	// api server routes
	r := h.Router(v.BuildTime, v.Commit, v.Release)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// this channel is for graceful shutdown:
	// if we receive an error, we can send it here to notify the server to be stopped
	shutdown := make(chan struct{}, 1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			shutdown <- struct{}{}
			log.Printf("%s %v", logHeader, err)
		}
	}()

	log.Println(u.Envelope(
		fmt.Sprintf("%s The druid is listen on http://0.0.0.0:%v/", logHeader, port)))

	//create new docker registry worker
	stopWorker := make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go w.NewWorker(time.Second*60, s.Config, s.DB, stopWorker, &wg)

	// wait for signal
	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			log.Println(u.Envelope(logHeader + " Got SIGINT..."))
		case syscall.SIGTERM:
			log.Println(u.Envelope(logHeader + " Got SIGTERM..."))
		}
	case <-shutdown:
		log.Println(u.Envelope(logHeader + " Got shutdown..."))
	}

	// killall
	log.Println(u.Envelope(logHeader + " The druid is shutting down..."))
	srv.Shutdown(context.Background())
	stopWorker <- true
	wg.Wait()
	log.Println(u.Envelope(logHeader + " Done"))
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/serhio83/druid/pkg/handlers"
	"github.com/serhio83/druid/pkg/utils"
	"github.com/serhio83/druid/pkg/version"
)

func main() {
	log.Println(utils.StringDecorator(fmt.Sprintf(
		"Starting the druid. Commit: %s, build time: %s, release: %s",
		version.Commit, version.BuildTime, version.Release,
	)))

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal(utils.StringDecorator("Port is not set"))
	}

	r := handlers.Router(version.BuildTime, version.Commit, version.Release)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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
			log.Printf("%v", err)
		}
	}()
	log.Println(utils.StringDecorator(
		fmt.Sprintf("The druid is listen on http://0.0.0.0:%v/", port)))

	select {
	case killSignal := <-interrupt:
		switch killSignal {
		case os.Interrupt:
			log.Println(utils.StringDecorator("Got SIGINT..."))
		case syscall.SIGTERM:
			log.Println(utils.StringDecorator("Got SIGTERM..."))
		}
	case <-shutdown:
		log.Println(utils.StringDecorator("Got an error..."))
	}

	log.Println(utils.StringDecorator("The druid is shutting down..."))
	srv.Shutdown(context.Background())
	log.Println(utils.StringDecorator("Done"))
}

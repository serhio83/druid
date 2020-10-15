package handlers

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
	u "github.com/serhio83/druid/pkg/utils"
)

const (
	logHeader = "[druid_api]"
)

//Router return new mux.Router
func Router(buildTime, commit, release string) *mux.Router {
	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Println(u.Envelope(logHeader + " Readyz probe is negative by default..."))
		time.Sleep(200 * time.Millisecond)
		isReady.Store(true)
		log.Println(u.Envelope(logHeader + " Readyz probe is positive."))
	}()

	r := mux.NewRouter()
	// r.HandleFunc("/", checkHeaders(mainpage))
	r.HandleFunc("/", home(buildTime, commit, release)).Methods("GET")
	r.HandleFunc("/healthz", healthz)
	r.HandleFunc("/readyz", readyz(isReady))
	return r
}

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	u "github.com/serhio83/druid/pkg/utils"
)

func checkHeaders(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rHost := strings.Split(r.Host, ":")[0]
		rAddr := strings.Split(r.RemoteAddr, ":")[0]

		// check if Content-Type: application/json used and fail if not
		ah := r.Header.Get("Content-Type")
		if ah != "application/json" {
			log.Println(u.Envelope(
				fmt.Sprintf("%s %s - %s [400] you should use Content-Type: application/json",
					logHeader,
					rAddr,
					rHost)))
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		// fail if zero content length
		if r.ContentLength == 0 {
			log.Println(u.Envelope(
				fmt.Sprintf("%s %s - %s [400] Invalid request payload",
					logHeader,
					rAddr,
					rHost)))
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		f(w, r)
	}
}

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/serhio83/druid/pkg/utils"
)

func checkHeaders(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rHost := strings.Split(r.Host, ":")[0]
		rAddr := strings.Split(r.RemoteAddr, ":")[0]

		// check if Content-Type: application/json used and fail if not
		ah := r.Header.Get("Content-Type")
		if ah != "application/json" {
			log.Println(utils.StringDecorator(
				fmt.Sprintf("%s - %s [400] you should use Content-Type: application/json",
					rAddr,
					rHost)))
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		// fail if zero content length
		if r.ContentLength == 0 {
			log.Println(utils.StringDecorator(
				fmt.Sprintf("%s - %s [400] Invalid request payload",
					rAddr,
					rHost)))
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		f(w, r)
	}
}

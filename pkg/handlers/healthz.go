package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	u "github.com/serhio83/druid/pkg/utils"
)

func healthz(w http.ResponseWriter, _ *http.Request) {
	okz := struct {
		Okz string `json:"okz"`
	}{"work fine"}
	body, err := json.Marshal(okz)
	if err != nil {
		log.Println(u.Envelope(fmt.Sprintf("Error marshaling json: %v", err)))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

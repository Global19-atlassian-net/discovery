package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// HealthHandler checks weather the discovery service is in a healthy status
func (h Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	token, err := h.setupToken(0)

	if err != nil || token == "" {
		log.Printf("health failed to setup token %v", err)
		http.Error(w, "health failed to setup token", 400)
		return
	}

	err = h.deleteToken(token)
	if err != nil {
		log.Printf("health failed to delete token %v", err)
		http.Error(w, "health failed to delete token", 400)
		return
	}

	fmt.Fprintf(w, "OK")
}

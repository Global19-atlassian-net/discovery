package handlers

import (
	"net/http"
)

// HomeHandler handles the home url of the discovery service
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		"https://github.com/quantum/discovery",
		http.StatusMovedPermanently,
	)
}

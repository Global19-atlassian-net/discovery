package handlers

import (
	"fmt"
	"net/http"
)

func (h *Handler) RobotsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "User-agent: *\nDisallow: /")
}

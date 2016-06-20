package badge

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/narqo/go-badge"
)

const (
	Title = "exago"
)

func Write(w http.ResponseWriter, title, val, color string) {
	if title == "" {
		title = Title
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	err := badge.Render(title, val, badge.Color(color), w)
	if err != nil {
		log.Error(err)
	}
}

func WriteError(w http.ResponseWriter, title string) {
	Write(w, title, "error", "lightgrey")
}

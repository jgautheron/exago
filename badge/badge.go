package badge

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/narqo/go-badge"
)

const (
	Title = "exago"
)

func Write(w http.ResponseWriter, val, color string) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	err := badge.Render(Title, val, badge.Color(color), w)
	if err != nil {
		log.Error(err)
	}
}

func WriteError(w http.ResponseWriter) {
	Write(w, "error", "lightgrey")
}

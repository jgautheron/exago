package badge

import (
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/narqo/go-badge"
)

const (
	Title = "exago"
)

func Write(w http.ResponseWriter, title, val string, score float64) {
	if title == "" {
		title = Title
	}

	color := "lightgrey"
	if score != -1 {
		color = "hsl(" + string(strconv.FormatFloat(score-20, 'f', 2, 64)) + ", 60%, 50%)"
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	err := badge.Render(title, val, badge.Color(color), w)
	if err != nil {
		log.Error(err)
	}
}

func WriteError(w http.ResponseWriter, title string) {
	Write(w, title, "error", -1)
}

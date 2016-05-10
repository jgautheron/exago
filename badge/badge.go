package badge

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const (
	title = "exago"
)

func Write(w http.ResponseWriter, val, color string) {
	// TODO: Should work directly instead on the raw SVG file
	resp, err := http.Get("https://img.shields.io/badge/" + title + "-" + val + "-" + color + ".svg")
	if err != nil {
		log.Error(err)
		WriteError(w)
		return
	}
	defer resp.Body.Close()
	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		WriteError(w)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(img)
}

func WriteError(w http.ResponseWriter) {
	Write(w, "error", "lightgrey")
}

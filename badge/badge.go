package badge

import (
	"io/ioutil"
	"net/http"
)

func Write(w http.ResponseWriter, name, val, color string) {
	// TODO: Should work directly instead on the SVG file
	resp, err := http.Get("https://img.shields.io/badge/" + name + "-" + val + "-" + color + ".svg")
	if err != nil {
		WriteError(w, name)
		return
	}
	defer resp.Body.Close()
	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		WriteError(w, name)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(img)
}

func WriteError(w http.ResponseWriter, name string) {
	Write(w, name, "error", "lightgrey")
}

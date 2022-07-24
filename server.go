package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/browser"
)

type ApiFilenamePayload struct {
	Filename string `json:"filename"`
}

func apiSaveFileHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("/api/save")

	var payload ApiFilenamePayload

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("%v\n", payload)
	go SaveArduinoDataToFile(payload.Filename)

	fmt.Fprintf(w, `{"ok":true}`)
}

type ApiPortPayload struct {
	Port string `json:"port"`
}

func apiSetPortHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("/api/port")

	var payload ApiPortPayload

	err := json.NewDecoder(req.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("%v\n", payload)
	err = ArduinoConnect(payload.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, `{"ok": true}`)
	}
}

//go:embed *.html *.js *.css
var content embed.FS

func main() {
	browser.OpenURL("http://localhost:8080")

	http.HandleFunc("/api/save", apiSaveFileHandler)
	http.HandleFunc("/api/port", apiSetPortHandler)

	http.Handle("/", http.FileServer(http.FS(content)))

	http.ListenAndServe(":8080", nil)
}

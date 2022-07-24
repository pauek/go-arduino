package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
)

var upgrader = websocket.Upgrader{} // use default options

type WebsocketCommand struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

func websocketHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("/ws")
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		reader := bytes.NewReader(msg)
		var cmd WebsocketCommand
		json.NewDecoder(reader).Decode(&cmd)
		log.Printf("cmd: %v\n", cmd)

		switch cmd.Cmd {
		case "saveFile":
			if len(cmd.Args) != 1 {
				ws.WriteMessage(msgType, []byte(
					`{"cmd": "saveFile", "ok": false, "error": "Args != 1"}`,
				))
				return
			}
			go SaveArduinoDataToFile(cmd.Args[0])
			ws.WriteMessage(msgType, []byte(
				`{"cmd": "saveFile", "ok": true}`,
			))

		case "setPort":
			if len(cmd.Args) != 1 {
				ws.WriteMessage(msgType, []byte(
					`{"cmd": "setPort", "ok": false, "error": "Args != 1"}`,
				))
			}
			err = ArduinoConnect(cmd.Args[0])
			if err != nil {
				ws.WriteMessage(msgType, []byte(fmt.Sprintf(
					`{"cmd": "setPort", "ok": false, "error": "%s"}`,
					err.Error(),
				)))
			} else {
				ws.WriteMessage(msgType, []byte(
					`{"cmd": "setPort", "ok": true}`,
				))
			}
		default:
			ws.WriteMessage(msgType, []byte(fmt.Sprintf(
				`{"cmd": "%s", "ok": false, "error": "%s"}`,
				cmd.Cmd,
				"Unrecognized command!",
			)))
		}
	}
}

//go:embed *.html *.js *.css
var content embed.FS

func main() {
	browser.OpenURL("http://localhost:8080")

	http.HandleFunc("/ws", websocketHandler)
	http.Handle("/", http.FileServer(http.FS(content)))

	http.ListenAndServe(":8080", nil)
}

package poker

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type playerServerWS struct {
	*websocket.Conn
}

func NewPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	// upgrades http connection to ws
	// then read and record user input
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to WebSocket, %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	// ReadMessage blocks on waiting for input
	_, msg, err := w.ReadMessage()

	if err != nil {
		log.Printf("error reading from websocket, %v\n", err)
	}

	return string(msg)
}

func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

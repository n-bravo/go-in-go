package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Upgrader websocket.Upgrader
}

var Manager *SessionManager = NewSessionManager()

func (wsh WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wsh.Upgrader.Upgrade(w, r, nil)	
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	m := HandshakeSessionMessage{}
	for {
		if err = c.ReadJSON(&m); err != nil {
			log.Printf("Error %s when reading handshake message from client", err)
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				return
			}
			continue
		}
		if m.SessionId == "" { //create new session
			if m.Size != 19 {
				msg := "error invalid board size"
				log.Println(msg)
				c.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
				continue
			}
			if m.Online {
				msg := "error no online support yet"
				log.Println(msg)
				c.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
				continue
			}
			log.Printf("Creating new session")
			Manager.NewSession(c, m.Size, m.Online)
			return
		}
	}
}

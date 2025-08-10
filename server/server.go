package server

import (
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	Upgrader websocket.Upgrader
	Origins  []string
}

var Manager *SessionManager = NewSessionManager()

func (wsh WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := wsh.Upgrader
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return slices.Contains(wsh.Origins, r.Header["Origin"][0])
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	m := HandshakeSessionMessage{}
	for {
		if err = c.ReadJSON(&m); err != nil {
			log.Printf("Error %s when reading handshake message from client", err)
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				c.Close()
				return
			}
			c.Close()
			return
		}
		if m.SessionId == "" { //create new session
			if m.Size != 19 && m.Size != 5 {
				msg := "error invalid board size"
				log.Println(msg)
				c.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
				c.Close()
				return
			}
			log.Printf("Creating new session")
			Manager.NewSession(c, m.Size, m.Online)
			return
		} else {
			if !Manager.OnlineSessionExists(m.SessionId) {
				msg := fmt.Sprintf("online session id %s not found", m.SessionId)
				log.Println(msg)
				c.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
				c.Close()
				return
			}
			Manager.JoinSession(m.SessionId, c)
			return
		}
	}
}

package server

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/n-bravo/go-in-go/game"
)

type session struct {
	id     string
	online bool
	conn   *websocket.Conn
	g      *game.GoGame
	m      *SessionManager
}

func newSession(c *websocket.Conn, n int, online bool, m *SessionManager) (*session, error) {
	g, err := game.NewGame(n)
	if err != nil {
		return nil, err
	}
	s := session{
		id:     uuid.NewString(),
		online: online,
		conn:   c,
		g:      g,
		m:      m,
	}
	if err = s.conn.WriteJSON(&NewSessionMessage{SessionId: s.id}); err != nil {
		return nil, fmt.Errorf("error when sending new session information to client: %s", err)
	}
	if online {
		
	} else {
		go s.offlineLoop()
	}
	return &s, nil
}

func (s *session) offlineLoop() {
	defer s.m.CloseSession(s)
	for {
		var err error
		var input PlayerInputMessage
		if err = s.conn.ReadJSON(&input); err != nil {
			msg := fmt.Sprintf("Error when reading input from client from session: %s", s.id)
			log.Println(msg)
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				return
			}
			continue
		}
		if input.CloseSession {
			log.Printf("Client request close session %s", s.id)
			return
		}
		if err = s.g.Play(input.X, input.Y, input.Black); err != nil {
			msg := fmt.Sprintf("Invalid request from client: %s", err)
			log.Println(msg)
			s.conn.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
			continue
		}
		s.conn.WriteJSON(&ResponseMessage{Code: 200, Message: ""})
	}
}

func (s *session) close() error {
	var err error
	log.Printf("Closing session %s", s.id)
	err = s.conn.Close()
	if err != nil {
		return err
	}
	err = s.g.Close()
	if err != nil {
		return err
	}
	return nil
}

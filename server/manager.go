package server

import (
	"log"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

type SessionManager struct {
	sessions map[*session]bool
}

func NewSessionManager() *SessionManager {
	uuid.EnableRandPool()
	return &SessionManager{
		sessions: make(map[*session]bool),
	}
}

func (m *SessionManager) NewSession(c *websocket.Conn, size int, online bool) {
	s, err := newSession(c, size, online, m)
	if err != nil {
		log.Fatal(err)
	}
	m.sessions[s] = true   //TODO: mutex
}

func (m *SessionManager) CloseSession(s *session) {
	if err := s.close(); err != nil {
		log.Fatal(err)
	}
	delete(m.sessions, s)
}

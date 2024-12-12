package server

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SessionManager struct {
	mu       sync.Mutex
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

	go func() {
		m.mu.Lock()
		m.sessions[s] = true
		m.mu.Unlock()
	}()
}

func (m *SessionManager) CloseSession(s *session) {
	if err := s.close(); err != nil {
		log.Fatal(err)
	}
	go func() {
		m.mu.Lock()
		delete(m.sessions, s)
		m.mu.Unlock()
	}()
}

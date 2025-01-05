package server

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SessionManager struct {
	mu       sync.Mutex
	sessions map[session]bool
}

func NewSessionManager() *SessionManager {
	uuid.EnableRandPool()
	return &SessionManager{
		sessions: make(map[session]bool),
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

func (m *SessionManager) OnlineSessionExists(id string) bool {
	for s := range m.sessions {
		if s.isOnline() && s.getId() == id {
			return true
		}
	}
	return false
}

func (m *SessionManager) JoinSession(id string, c *websocket.Conn) {
	for s := range m.sessions {
		if s.getId() == id {
			s.addPlayer(c)
			return
		}
	}
}

func (m *SessionManager) CloseSession(s session) {
	go func(){
		m.mu.Lock()
		delete(m.sessions, s)
		m.mu.Unlock()
	}()
}

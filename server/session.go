package server

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/n-bravo/go-in-go/game"
)

type session interface {
	getId() string
	mainLoop()
	addPlayer(c *websocket.Conn) error
	close(con *websocket.Conn) error
}

type offlineSession struct {
	id   string
	conn *websocket.Conn
	g    *game.GoGame
	m    *SessionManager
}

type onlineSession struct {
	id         string
	mu         sync.Mutex
	con1, con2 *websocket.Conn //con1 is always black, con2 is always white
	g          *game.GoGame
	m          *SessionManager
}

func newSession(c *websocket.Conn, n int, online bool, m *SessionManager) (session, error) {
	g, err := game.NewGame(n)
	if err != nil {
		return nil, err
	}
	if online {
		s := &onlineSession{
			id:   uuid.NewString(),
			con1: c,
			con2: nil,
			g:    g,
			m:    m,
		}
		if err = s.con1.WriteJSON(&NewSessionMessage{SessionId: s.id}); err != nil {
			return nil, fmt.Errorf("error when sending new session information to client: %s", err)
		}
		go s.mainLoop()
		return s, nil
	} else {
		s := &offlineSession{
			id:   uuid.NewString(),
			conn: c,
			g:    g,
			m:    m,
		}
		if err = s.conn.WriteJSON(&NewSessionMessage{SessionId: s.id}); err != nil {
			return nil, fmt.Errorf("error when sending new session information to client: %s", err)
		}
		go s.mainLoop()
		return s, nil
	}
}

func (s *offlineSession) getId() string {
	return s.id
}

func (s *onlineSession) getId() string {
	return s.id
}

func (s *offlineSession) addPlayer(c *websocket.Conn) error {
	return nil
}

func (s *onlineSession) addPlayer(c *websocket.Conn) error {
	defer s.mu.Unlock()
	s.mu.Lock()
	if s.con2 != nil {
		return fmt.Errorf("error session %s is already full", s.id)
	}
	s.con2 = c
	s.mu.Unlock()
	log.Printf("Player 2 joined to session %s", s.id)
	c.WriteJSON(&ResponseMessage{Code: 200})
	return nil
}

func (s *offlineSession) mainLoop() {
	defer s.m.CloseSession(s, s.conn)
	for {
		var err error
		var input OffilePlayerInputMessage
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

func (s *onlineSession) mainLoop() {
	go s.onlinePlayerLoop(true)
	for {
		//wait for con2 to connect
		s.mu.Lock()
		if s.con2 != nil {
			break
		}
		s.mu.Unlock()
	}
	log.Printf("Session %s with all their players ready to play", s.id)
	go s.onlinePlayerLoop(false)
}

func (s *onlineSession) onlinePlayerLoop(black bool) {
	var con *websocket.Conn
	var pname byte
	if black {
		con = s.con1
		pname = '1'
	} else {
		con = s.con2
		pname = '2'
	}
	defer func() {
		s.m.CloseSession(s, con)
	}()
	for {
		var err error
		var input OnlinePlayerInputMessage
		if err = con.ReadJSON(&input); err != nil {
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				return
			}
			msg := fmt.Sprintf("Error when reading input from client %s [%s]: %v", string(pname), s.id, err)
			log.Println(msg)
			continue
		}
		if input.CloseConn {
			log.Printf("Client %s [%s] request close session", string(pname), s.id)
			return
		}
		s.mu.Lock()
		if s.con1 == nil || s.con2 == nil {
			s.mu.Unlock()
			msg := fmt.Sprintf("error in session %s: all players are not connected", s.id)
			log.Println(msg)
			con.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
			continue
		}
		if err = s.g.Play(input.X, input.Y, black); err != nil {
			s.mu.Unlock()
			msg := fmt.Sprintf("Invalid request from client %s [%s]: %s", string(pname), s.id, err)
			log.Println(msg)
			con.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
			s.mu.Unlock()
			continue
		}
		con.WriteJSON(&ResponseMessage{Code: 200, Message: ""})
	}
}

func (s *offlineSession) close(con *websocket.Conn) error {
	var err error
	log.Printf("Closing session %s", s.id)
	err = con.Close()
	if err != nil {
		return err
	}
	err = s.g.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *onlineSession) close(con *websocket.Conn) error {
	var err error
	if s.con1 == con {
		log.Printf("Closing client 1 of session %s", s.id)
		err = s.con1.Close()
		if err != nil {
			return err
		}
		s.con1 = nil
	} else {
		log.Printf("Closing client 2 of session %s", s.id)
		err = s.con2.Close()
		if err != nil {
			return err
		}
		s.con2 = nil
	}
	if s.con1 == nil && s.con2 == nil {
		log.Printf("Closing session %s", s.id)
		err = s.g.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

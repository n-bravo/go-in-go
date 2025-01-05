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
	addPlayer(c *websocket.Conn)
    isOnline() bool
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
	mu         sync.Mutex      //only used to block the access to the board game g
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
		if err = s.con1.WriteJSON(&NewSessionMessage{SessionId: s.id, Online: true, BlackSide: true}); err != nil {
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
		if err = s.conn.WriteJSON(&NewSessionMessage{SessionId: s.id, Online: false, BlackSide: true}); err != nil {
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

func (s *offlineSession) isOnline() bool {
    return false
}

func (s *onlineSession) isOnline() bool {
    return true
}

func (s *offlineSession) addPlayer(c *websocket.Conn) {
}

func (s *onlineSession) addPlayer(c *websocket.Conn) {
	if s.con1 != nil && s.con2 != nil {
		msg := fmt.Sprintf("error session %s is already full", s.id)
		log.Print(msg)
		c.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
		c.Close()
		return
	}
	if s.con1 == nil {
		log.Printf("Player 1 joined to session %s", s.id)
		go s.onlinePlayerLoop(true)
		s.con1 = c
        c.WriteJSON(&NewSessionMessage{SessionId: s.id, Online: true, BlackSide: true, BStatus: s.g.String()})
	} else {
		log.Printf("Player 2 joined to session %s", s.id)
		go s.onlinePlayerLoop(false)
		s.con2 = c
        c.WriteJSON(&NewSessionMessage{SessionId: s.id, Online: true, BlackSide: false, BStatus: s.g.String()})
	}
}

func (s *offlineSession) mainLoop() {
	defer s.m.CloseSession(s)
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
		s.close(con)
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
		if s.con1 == nil || s.con2 == nil {
			msg := fmt.Sprintf("error in session %s: all players are not connected", s.id)
			log.Println(msg)
			con.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
			continue
		}
		s.mu.Lock()
		if err = s.g.Play(input.X, input.Y, black); err != nil {
			s.mu.Unlock()
			msg := fmt.Sprintf("Invalid request from client %s [%s]: %s", string(pname), s.id, err)
			log.Println(msg)
			con.WriteJSON(&ResponseMessage{Code: 401, Message: msg})
			continue
		}
        status := s.g.String()
		s.mu.Unlock()
        con.WriteJSON(&ResponseMessage{Code: 200, Message: "", BStatus: status})
        if con == s.con1 {
            s.con2.WriteJSON(&ResponseMessage{Code: 200, Message: "", BStatus: status})
        } else {
            s.con1.WriteJSON(&ResponseMessage{Code: 200, Message: "", BStatus: status})
        }
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
		s.m.CloseSession(s)
	}
	return nil
}

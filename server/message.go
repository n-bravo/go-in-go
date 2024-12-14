package server

type HandshakeSessionMessage struct {
	SessionId string `json:"sessionId"`
	Size      int    `json:"size"`
	Online    bool   `json:"online"`
}

type OffilePlayerInputMessage struct {
	X            int  `json:"x"`
	Y            int  `json:"y"`
	Black        bool `json:"black"`
	CloseSession bool `json:"closeSession"`
}

type OnlinePlayerInputMessage struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	CloseConn bool `json:"closeConn"`
}

type NewSessionMessage struct {
	SessionId string `json:"sessionId"`
}

type ResponseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

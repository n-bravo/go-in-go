package server

// Message formats supported for websocket communications

// Initial from client to server to open a websocket connection
type HandshakeSessionMessage struct {
	SessionId string `json:"sessionId"` // ID of the session to join. Empty string if want to create a new session
	Size      int    `json:"size"`      // Board size of the new session. Only 5 and 19 supported currently. Ignored if SessionId is not empty.
	Online    bool   `json:"online"`    // 'true' if want to create a new online session. 'false' otherwise. Ignored if SessionId is not empty.
}

// Response from server to client after a HandshakeSessionMessage is process.
// Creating a new session or joining the client to an existing one.
//
// bStatus is a string representation of the board when joining an existing session.
// It represents each intersection in the board, using a \n character to separete each row.
// * = empty intersection
// B = intersection taken by black stones
// W = intersection taken by white stones
type NewSessionResponseMessage struct {
	SessionId string `json:"sessionId"` // ID of the new session or the session joined.
	Online    bool   `json:"online"`    // true if the session is online. false otherwise.
	BlackSide bool   `json:"blackSide"` // true if the client is assigned to black side. false if assigned to white side.
	BStatus   string `json:"bStatus"`   // Board status when creating or joining the session.
}

// User movement action message for an offline match.
// Even though the match is offline (both players using the same client), each movement is sent to
// the server to verify turn order, captures and valid movements.
type OffilePlayerInputMessage struct {
	X            int  `json:"x"`            // X position of the movement
	Y            int  `json:"y"`            // Y position of the movement
	Black        bool `json:"black"`        // true if the movement correspond to black side, false otherwise
	CloseSession bool `json:"closeSession"` // true if want to close the connection, finishing the session. Omit or false otherwise.
}

// User movement action message for an online match.
// The side is assigned according to the connection order. The session creator is black side, and the client joining after is white side.
type OnlinePlayerInputMessage struct {
	X         int  `json:"x"`         // X position of the movement
	Y         int  `json:"y"`         // Y position of the movement
	CloseConn bool `json:"closeConn"` // true if want to close the connection. Omit or false otherwise. The session will be still alive as long one client is connected.

}

// Response from server to client after a new movement from the client
type ResponseMessage struct {
	Code    int    `json:"code"`    // HTTP convention (for easy understanding). 200 is a correct move. 401 is a forbidden move (either by wrong turn order or invalid position)
	Message string `json:"message"` // In case Code is not 200, the server will provide a message to explaing why.
	BStatus string `json:"bStatus"` // Board status after a valid client movement. Same format as NewSessionResponseMessage.BStatus
}

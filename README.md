# go-in-go

As a little fun project to learn Golang, **Go in Go** is a websocket server to handle online [Go](https://en.wikipedia.org/wiki/Go_(game)) matches.
The core features of the server are:

- Authenticate users (TODO)
- Create online matches through websocket messages.
- Create offline matches. (TODO)
- Handle matches actions:
    - User movements
    - Turn order
    - Territory captures
    - Verify valid user actions
    - User scores (TODO)
- Persist matches for long time pauses or unexpected disconnections. (TODO)

## Requirements

- Go 1.23.3

## How to run

The server currently does not need external dependencies, so just type:

```bash
go mod tidy
go run .
```

## Server messages

The websocket server comunicates with the client with a series of messages in JSON format,
defined in `server/message.go`:

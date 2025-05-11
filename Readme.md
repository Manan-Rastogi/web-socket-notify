## What is a WebSocket?
A WebSocket is a persistent, full-duplex connection between client and server over a single TCP connection — unlike HTTP which is request/response based.
-   You send a message from client → server without making a new request.
-   Server can push data to client at any time, without being asked.


## How does it start? (The WebSocket handshake)
The connection starts as a normal HTTP request:

```
GET /ws?deviceID=abc123 HTTP/1.1
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: ...
```
Then the server responds with:

```
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
```
This “Upgrade” handshake is what gorilla/websocket.Upgrader does for us under the hood.

## What does gorilla/websocket do?
It:

- Upgrades the net/http request to a WebSocket connection.

- Provides helper functions:
    -   ReadMessage() → reads one complete message from the client.
    -   WriteMessage() → sends a message to the client.

-   Abstracts frames, opcodes, and low-level socket operations.

## Why defer conn.Close()?
To gracefully close the connection when the function ends:

```go
defer conn.Close()
```
This ensures:

-   Connection is closed if client disconnects or server throws an error.
-   No memory leaks.

## What is mt in ReadMessage()?
```go
mt, msg, err := conn.ReadMessage()
```
mt = message type, like:

-   1 → text

-   2 → binary

-   9 → ping

-   10 → pong

You echo the message type back to preserve its meaning.



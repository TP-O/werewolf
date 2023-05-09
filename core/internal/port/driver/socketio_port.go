package driver

import socketio "github.com/googollee/go-socket.io"

type SocketIoPort interface {
	Connect(client socketio.Conn) error
	Disconnect(client socketio.Conn, reason string)
	HandleError(client socketio.Conn, err error)
	SyncPosition(client socketio.Conn, msg string)
}

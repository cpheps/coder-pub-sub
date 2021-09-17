package websocket

import (
	"context"
	"io"
	"net/http"

	"github.com/stretchr/testify/mock"
)

var _ (WebsocketConnection) = (*MockWebsocketConnection)(nil)

// MockWebsocketConnection represents a mock WebsocketConnection
type MockWebsocketConnection struct {
	mock.Mock
}

func (m *MockWebsocketConnection) NextWriter(messageType MessageType) (io.WriteCloser, error) {
	args := m.Called(messageType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(io.WriteCloser), args.Error(1)
}

func (m *MockWebsocketConnection) Close() error {
	args := m.Called()
	return args.Error(0)
}

var _ (io.WriteCloser) = (*MockWriteCloser)(nil)

// MockWriteCloser represents a mock io.WriteCloser
type MockWriteCloser struct {
	mock.Mock
}

func (m *MockWriteCloser) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockWriteCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

var _ (Broadcaster) = (*MockBroadcaster)(nil)

// MockBroadcaster represents a mock Broadcaster
type MockBroadcaster struct {
	mock.Mock
}

// RegisterConnection registers a connection with the Broadcaster
func (m *MockBroadcaster) RegisterConnection(conn WebsocketConnection) {
	m.Called(conn)
}

// Broadcast sends the bytes of messageType to all websockets.
// Returns and error if a single send fails
func (m *MockBroadcaster) Broadcast(ctx context.Context, messageType MessageType, msg []byte) error {
	args := m.Called(ctx, messageType, msg)
	return args.Error(0)
}

// CloseConnections closes all registered connections
func (m *MockBroadcaster) CloseConnections() {
	m.Called()
}

var _ (Upgrader) = (*MockUpgrader)(nil)

// MockUpgrader represents a mock Upgrader
type MockUpgrader struct {
	mock.Mock
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
func (m *MockUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (WebsocketConnection, error) {
	args := m.Called(w, r, responseHeader)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(WebsocketConnection), args.Error(1)
}

package websocket

import (
	"io"

	"github.com/stretchr/testify/mock"
)

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

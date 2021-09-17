package websocket

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewCacheBroadcaster(t *testing.T) {
	testCases := []struct {
		desc        string
		input       int
		expected    *CacheBroadcaster
		expectedErr error
	}{
		{
			desc:        "Invalid concurrency value",
			input:       -1,
			expected:    nil,
			expectedErr: errors.New("concurrency must be greater than 0"),
		},
		{
			desc:  "Valid create",
			input: 2,
			expected: &CacheBroadcaster{
				conns:       make([]WebsocketConnection, 0),
				concurrency: 2,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err := NewCacheBroadcaster(tc.input)

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedErr.Error())
			}
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_CacheBroadCaster_RegisterConnection(t *testing.T) {
	broadcaster, err := NewCacheBroadcaster(1)
	assert.NoError(t, err)

	// Create an empty connection since we don't care if it works
	connection := &GorillaConn{}

	broadcaster.RegisterConnection(connection)

	assert.Len(t, broadcaster.conns, 1)
	assert.Equal(t, connection, broadcaster.conns[0])
}

func Test_CacheBroadCaster_Broadcast(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(t *testing.T)
	}{
		{
			desc: "Conn NextWriter Failure",
			testFunc: func(t *testing.T) {
				messageType := TextMessage
				msg := []byte("hi")
				expectedErr := errors.New("bad stuff")

				mockConn := &MockWebsocketConnection{}
				mockConn.On("NextWriter", messageType).Return(nil, expectedErr)

				broadcaster, err := NewCacheBroadcaster(1)
				assert.NoError(t, err)

				broadcaster.RegisterConnection(mockConn)

				err = broadcaster.Broadcast(context.Background(), messageType, msg)

				assert.ErrorIs(t, err, expectedErr)
			},
		},
		{
			desc: "Write Failure",
			testFunc: func(t *testing.T) {
				messageType := TextMessage
				msg := []byte("hi")
				expectedErr := errors.New("bad stuff")

				mockWriter := &MockWriteCloser{}
				mockWriter.On("Write", msg).Return(0, expectedErr)

				mockConn := &MockWebsocketConnection{}
				mockConn.On("NextWriter", messageType).Return(mockWriter, nil)

				broadcaster, err := NewCacheBroadcaster(1)
				assert.NoError(t, err)

				broadcaster.RegisterConnection(mockConn)

				err = broadcaster.Broadcast(context.Background(), messageType, msg)

				assert.ErrorIs(t, err, expectedErr)
			},
		},
		{
			desc: "Write Success, close failure",
			testFunc: func(t *testing.T) {
				messageType := TextMessage
				msg := []byte("hi")
				expectedErr := errors.New("bad stuff")

				mockWriter := &MockWriteCloser{}
				mockWriter.On("Write", msg).Return(len(msg), nil)
				mockWriter.On("Close").Return(expectedErr)

				mockConn := &MockWebsocketConnection{}
				mockConn.On("NextWriter", messageType).Return(mockWriter, nil)

				broadcaster, err := NewCacheBroadcaster(1)
				assert.NoError(t, err)

				broadcaster.RegisterConnection(mockConn)

				err = broadcaster.Broadcast(context.Background(), messageType, msg)

				assert.ErrorIs(t, err, expectedErr)
			},
		},
		{
			desc: "Write Success",
			testFunc: func(t *testing.T) {
				messageType := TextMessage
				msg := []byte("hi")

				mockWriter := &MockWriteCloser{}
				mockWriter.On("Write", msg).Return(len(msg), nil)
				mockWriter.On("Close").Return(nil)

				mockConn := &MockWebsocketConnection{}
				mockConn.On("NextWriter", messageType).Return(mockWriter, nil)

				broadcaster, err := NewCacheBroadcaster(1)
				assert.NoError(t, err)

				broadcaster.RegisterConnection(mockConn)

				err = broadcaster.Broadcast(context.Background(), messageType, msg)

				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

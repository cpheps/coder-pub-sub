package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cpheps/coder-pub-sub/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_PubSubServer_New_Error(t *testing.T) {
	pubsubServer, err := New("", -1)
	assert.Error(t, err)
	assert.Nil(t, pubsubServer)
}

func Test_PubSubServer_Publish(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Broadcast fails",
			testFunc: func(t *testing.T) {
				expectedCode := http.StatusInternalServerError
				expectedResp := errorResponse{
					Message: "Internal Error",
				}

				message := []byte("hi")

				mockBroadcaster := &websocket.MockBroadcaster{}
				mockBroadcaster.On("Broadcast", mock.Anything, websocket.TextMessage, message).Return(errors.New("bad thing"))

				pubsubServer := &PubSubServer{
					broadcaster: mockBroadcaster,
				}

				req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/plubish", bytes.NewReader(message)).WithContext(context.Background())
				w := httptest.NewRecorder()

				pubsubServer.Publish(w, req)

				defer w.Result().Body.Close()
				data, err := io.ReadAll(w.Result().Body)
				assert.NoError(t, err)

				var resp errorResponse
				err = json.Unmarshal(data, &resp)
				assert.NoError(t, err)

				assert.Equal(t, expectedCode, w.Result().StatusCode)
				assert.Equal(t, expectedResp, resp)
			},
		},
		{
			desc: "Broadcast Success",
			testFunc: func(t *testing.T) {
				expectedCode := http.StatusNoContent

				message := []byte("hi")

				mockBroadcaster := &websocket.MockBroadcaster{}
				mockBroadcaster.On("Broadcast", mock.Anything, websocket.TextMessage, message).Return(nil)

				pubsubServer := &PubSubServer{
					broadcaster: mockBroadcaster,
				}

				req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/plubish", bytes.NewReader(message)).WithContext(context.Background())
				w := httptest.NewRecorder()

				pubsubServer.Publish(w, req)

				defer w.Result().Body.Close()

				assert.Equal(t, expectedCode, w.Result().StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}

}

func Test_PubSubServer_RegisterSubscriber(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "Upgrade Fails",
			testFunc: func(t *testing.T) {
				expectedCode := http.StatusInternalServerError
				expectedResp := errorResponse{
					Message: "Internal Error",
				}

				req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/subscribe", http.NoBody).WithContext(context.Background())
				w := httptest.NewRecorder()

				mockUpgrader := &websocket.MockUpgrader{}
				mockUpgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("bad thing"))

				doneChan := make(chan struct{})

				pubsubServer := &PubSubServer{
					doneChan: doneChan,
					upgrader: mockUpgrader,
				}

				// Ensure test doesn't block
				close(doneChan)

				pubsubServer.RegisterSubscriber(w, req)

				defer w.Result().Body.Close()
				data, err := io.ReadAll(w.Result().Body)
				assert.NoError(t, err)

				var resp errorResponse
				err = json.Unmarshal(data, &resp)
				assert.NoError(t, err)

				assert.Equal(t, expectedCode, w.Result().StatusCode)
				assert.Equal(t, expectedResp, resp)
			},
		},
		{
			desc: "Success",
			testFunc: func(t *testing.T) {
				expectedCode := http.StatusOK

				req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/subscribe", http.NoBody).WithContext(context.Background())
				w := httptest.NewRecorder()

				mockWebsocket := &websocket.MockWebsocketConnection{}

				mockUpgrader := &websocket.MockUpgrader{}
				mockUpgrader.On("Upgrade", mock.Anything, mock.Anything, mock.Anything).Return(mockWebsocket, nil)

				mockBroadcaster := &websocket.MockBroadcaster{}
				mockBroadcaster.On("RegisterConnection", mockWebsocket)

				doneChan := make(chan struct{})

				pubsubServer := &PubSubServer{
					doneChan:    doneChan,
					upgrader:    mockUpgrader,
					broadcaster: mockBroadcaster,
				}

				// Ensure test doesn't block
				close(doneChan)

				pubsubServer.RegisterSubscriber(w, req)

				defer w.Result().Body.Close()

				assert.Equal(t, expectedCode, w.Result().StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}

}

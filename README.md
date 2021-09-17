# coder-pub-sub

## Requirements

Create an HTTP server that implements an in-memory, realtime PubSub system.

### Endpoints

| Endpoint | Description |
| :-- | :-- |
| Register Subscriber | Creates a persistent HTTP connection with the server. One way communication to receive messages |
| Publish | Sends a messages to the server that is broadcasted to all subscribers |

### Goals
- [X] Include at least one automated test.
- [X] Adhere to Go best practices.
- [X] No race conditions.


## Server

The PubSub server is an HTTP server that has the following endpoints:

| Path | Method | Paylog | Description |
| :--: | :--: | :--: | :-- |
| /subscribe | GET | None | Registers a subscriber with the server. This is a websocket connection so websocket headers are needed (show in example below) to establish a valid connection. The connection will persist until the server stops |
| /publish | POST | Any valid string | Takes the Test payload and forwards it onto all subscribers |

### Getting Started

The server can be run by running

```sh
go run main.go
```

This will start a server listening on port `8080`.

To register a subscriber connection run the following curl command:

```sh
curl --include \
     --no-buffer \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Host: example.com:80" \
     --header "Origin: http://example.com:80" \
     --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
     --header "Sec-WebSocket-Version: 13" \
     http://localhost:8080/subscribe
```

**Note** the value of `Sec-WebSocket-Key` is just a base64 encoded string to assist in the handshake and for demo purposes can be left alone.

Once you have your subscribers setup run the following curl command with the **text** payload of your choice to see it publish to the subscribers

```sh
curl -X POST -d "my cool payload" http://localhost:8080/publish
```

To stop the demo just Ctrl+C the server and everything will clean up.

## Things I would have added if real

Below are a list of things I would have done if this were to be a real service:

- Added a configuration for the server
- Handle dropped connection by client so we don't try to send down a broken pipeline
- Parse content type of publish to handle more payload types
- Added integration test to test the server as a standalone entity
- Added better API documentation
- Better logging (use a library that's more robust than the stdlib `log`)
- This was out of spec but adding traditional PubSub behavior like Topics, acking a message, storing messages

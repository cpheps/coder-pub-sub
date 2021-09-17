# coder-pub-sub

## Requirements

Create an HTTP server that implements an in-memory, realtime PubSub system.

### Endpoints

| Endpoint | Description |
| :--: | :--: |
| Register Subscriber | Creates a persistent HTTP connection with the server. One way communication to receive messages |
| Publish | Sends a messages to the server that is broadcasted to all subscribers |

### Goals
- [ ] Include at least one automated test.
- [ ] Adhere to Go best practices.
- [ ] No race conditions.

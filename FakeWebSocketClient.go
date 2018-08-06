package main

type FakeWebSocketClient struct {
	SendMessageArgs [][]byte
	CloseCalled     bool
}

func NewFakeWebSocketClient() *FakeWebSocketClient {
	return &FakeWebSocketClient{}
}

func (c *FakeWebSocketClient) MessageStream() <-chan WebSocketMessage {
	return make(chan WebSocketMessage)
}

func (c *FakeWebSocketClient) SendMessage(msg []byte) error {
	c.SendMessageArgs = append(c.SendMessageArgs, msg)
	return nil
}

func (c *FakeWebSocketClient) Close() {
	c.CloseCalled = true
}

package main

import "errors"

type FakeBrokenWebSocketClient struct {
	CloseCalled bool
}

func NewFakeBrokenWebSocketClient() *FakeBrokenWebSocketClient {
	return &FakeBrokenWebSocketClient{}
}

func (c *FakeBrokenWebSocketClient) MessageStream() <-chan WebSocketMessage {
	return make(chan WebSocketMessage)
}

func (c *FakeBrokenWebSocketClient) SendMessage(msg []byte) error {
	return errors.New("")
}

func (c *FakeBrokenWebSocketClient) Close() {
	c.CloseCalled = true
}

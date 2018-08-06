package main_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "private-chat"
)

var _ = Describe("WebSocketClientHub", func() {
	Context("Initializing", func() {
		It("Can be initalized and start", func() {
			go NewWebSocketClientHub().Start()
		})
	})

	Context("Registering and publishing messages", func() {
		It("Can send a messages to registered client", func() {
			webSocketClientHub := NewWebSocketClientHub()
			go webSocketClientHub.Start()

			webSocketClient := NewFakeWebSocketClient()
			webSocketClientHub.Register(webSocketClient)

			webSocketClientHub.Publish([]byte("message one"))
			webSocketClientHub.Publish([]byte("message two"))

			Expect(len(webSocketClient.SendMessageArgs)).To(Equal(2))
			Expect(webSocketClient.SendMessageArgs).To(ContainElement([]byte("message one")))
			Expect(webSocketClient.SendMessageArgs).To(ContainElement([]byte("message two")))

		})

		It("Can send a messages to multiple registered clients", func() {
			webSocketClientHub := NewWebSocketClientHub()
			go webSocketClientHub.Start()

			webSocketClientA := NewFakeWebSocketClient()
			webSocketClientB := NewFakeWebSocketClient()

			webSocketClientHub.Register(webSocketClientA)
			webSocketClientHub.Register(webSocketClientB)

			webSocketClientHub.Publish([]byte("message one"))

			Expect(len(webSocketClientA.SendMessageArgs)).To(Equal(1))
			Expect(len(webSocketClientB.SendMessageArgs)).To(Equal(1))
		})
	})

	Context("Unregistering and publishing messages", func() {

		It("Does not send messages to unregistered clients", func() {
			webSocketClientHub := NewWebSocketClientHub()
			go webSocketClientHub.Start()

			webSocketClientA := NewFakeWebSocketClient()
			webSocketClientB := NewFakeWebSocketClient()

			webSocketClientHub.Register(webSocketClientA)
			webSocketClientHub.Register(webSocketClientB)
			webSocketClientHub.Publish([]byte("message one"))

			webSocketClientHub.Unregister(webSocketClientB)
			webSocketClientHub.Publish([]byte("message one"))

			Expect(len(webSocketClientA.SendMessageArgs)).To(Equal(2))
			Expect(len(webSocketClientB.SendMessageArgs)).To(Equal(1))

		})

		It("Closes clients when they are unregistered", func() {
			webSocketClientHub := NewWebSocketClientHub()
			go webSocketClientHub.Start()

			webSocketClient := NewFakeWebSocketClient()

			webSocketClientHub.Register(webSocketClient)
			webSocketClientHub.Unregister(webSocketClient)
			fmt.Println("Test finished 111")

			Expect(webSocketClient.CloseCalled).To(Equal(true))
		})

		It("Unregisters broken clients", func() {
			webSocketClientHub := NewWebSocketClientHub()
			go webSocketClientHub.Start()

			webSocketClient := NewFakeBrokenWebSocketClient()
			webSocketClientHub.Register(webSocketClient)

			webSocketClientHub.Publish([]byte("message one"))
			fmt.Println("Test finished 222")
			Expect(webSocketClient.CloseCalled).To(Equal(true))
		})
	})
})

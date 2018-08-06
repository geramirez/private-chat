package main_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "private-chat"
)

var _ = Describe("WebSocketServiceHub", func() {
	Context("Initializing", func() {
		It("Can be initalized and start", func() {
			go NewWebSocketServiceHub().Start()
		})
	})

	Context("Registering and publishing messages", func() {
		It("Can send a messages to registered client", func() {
			webSocketServiceHub := NewWebSocketServiceHub()
			go webSocketServiceHub.Start()

			webSocketService := NewFakeWebSocketService()
			webSocketServiceHub.Register(webSocketService)

			webSocketServiceHub.Publish([]byte("message one"))
			webSocketServiceHub.Publish([]byte("message two"))

			Expect(len(webSocketService.SendMessageArgs)).To(Equal(2))
			Expect(webSocketService.SendMessageArgs).To(ContainElement([]byte("message one")))
			Expect(webSocketService.SendMessageArgs).To(ContainElement([]byte("message two")))

		})

		It("Can send a messages to multiple registered clients", func() {
			webSocketServiceHub := NewWebSocketServiceHub()
			go webSocketServiceHub.Start()

			webSocketServiceA := NewFakeWebSocketService()
			webSocketServiceB := NewFakeWebSocketService()

			webSocketServiceHub.Register(webSocketServiceA)
			webSocketServiceHub.Register(webSocketServiceB)

			webSocketServiceHub.Publish([]byte("message one"))

			Expect(len(webSocketServiceA.SendMessageArgs)).To(Equal(1))
			Expect(len(webSocketServiceB.SendMessageArgs)).To(Equal(1))
		})
	})

	Context("Unregistering and publishing messages", func() {

		It("Does not send messages to unregistered clients", func() {
			webSocketServiceHub := NewWebSocketServiceHub()
			go webSocketServiceHub.Start()

			webSocketServiceA := NewFakeWebSocketService()
			webSocketServiceB := NewFakeWebSocketService()

			webSocketServiceHub.Register(webSocketServiceA)
			webSocketServiceHub.Register(webSocketServiceB)
			webSocketServiceHub.Publish([]byte("message one"))

			webSocketServiceHub.Unregister(webSocketServiceB)
			webSocketServiceHub.Publish([]byte("message one"))

			Expect(len(webSocketServiceA.SendMessageArgs)).To(Equal(2))
			Expect(len(webSocketServiceB.SendMessageArgs)).To(Equal(1))

		})

		It("Closes clients when they are unregistered", func() {
			webSocketServiceHub := NewWebSocketServiceHub()
			go webSocketServiceHub.Start()

			webSocketService := NewFakeWebSocketService()

			webSocketServiceHub.Register(webSocketService)
			webSocketServiceHub.Unregister(webSocketService)
			fmt.Println("Test finished 111")

			Expect(webSocketService.CloseCalled).To(Equal(true))
		})

		It("Unregisters broken clients", func() {
			webSocketServiceHub := NewWebSocketServiceHub()
			go webSocketServiceHub.Start()

			webSocketService := NewFakeWebBrokenSocketService()
			webSocketServiceHub.Register(webSocketService)

			webSocketServiceHub.Publish([]byte("message one"))
			fmt.Println("Test finished 222")
			Expect(webSocketService.CloseCalled).To(Equal(true))
		})
	})
})

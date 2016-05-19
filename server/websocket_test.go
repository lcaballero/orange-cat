package server

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Websocket", func() {
	Describe("#NewWebsocket()", func() {
		It("should return a new Websocket object.", func() {
			sock := NewWebsocket("readme.md")
			Expect(sock).NotTo(BeNil())
		})
	})
})

package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPrivateChat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PrivateChat Suite")
}

package voldialogs

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVoldialogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Volume Dialogs Suite")
}

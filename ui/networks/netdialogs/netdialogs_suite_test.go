package netdialogs

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNetdialogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Network Dialogs Suite")
}

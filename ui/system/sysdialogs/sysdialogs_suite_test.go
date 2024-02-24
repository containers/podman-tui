package sysdialogs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSysdialogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sysdialogs Suite")
}

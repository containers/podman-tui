package dialogs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDialogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dialogs Suite")
}

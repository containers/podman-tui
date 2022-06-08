package imgdialogs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestImgdialogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Image dialogs Suite")
}

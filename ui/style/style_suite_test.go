package style_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStyle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UI Style Suite")
}

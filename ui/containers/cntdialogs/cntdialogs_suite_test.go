package cntdialogs_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCntdialogs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Containers Dialogs Suite")
}

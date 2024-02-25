package utils

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("home", func() {

	It("user home dir", func() {
		homedirEnv := os.Getenv("HOME")
		homedir, err := UserHomeDir()
		Expect(err).To(BeNil())
		Expect(homedir).To(Equal(homedirEnv))
	})
})

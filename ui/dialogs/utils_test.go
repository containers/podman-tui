package dialogs

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("command dialog", Ordered, func() {

	var sample1 string = "test\nabcdefghi\ntest"
	var sample2 string = "abcdefghi\ntest\ntest"
	var sample3 string = "test\ntest\nabcdefghi"
	var wants int = 9

	Describe("get message width", func() {
		It("sample1", func() {
			Expect(getMessageWidth(sample1)).To(Equal(wants))
		})
		It("sample2", func() {
			Expect(getMessageWidth(sample2)).To(Equal(wants))
		})
		It("sample3", func() {
			Expect(getMessageWidth(sample3)).To(Equal(wants))
		})
	})

})

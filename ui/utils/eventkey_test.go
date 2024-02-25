package utils

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("event key", func() {

	It("StringToEventKey", func() {
		input1 := "a"
		var input1events []*tcell.EventKey
		for i := 0; i < len(input1); i++ {
			input1events = append(input1events, tcell.NewEventKey(256, rune(input1[i]), tcell.ModNone))
		}

		Expect(StringToEventKey(input1)[0].Key()).To(Equal(input1events[0].Key()))
	})
})

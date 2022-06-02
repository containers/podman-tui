package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("prgbar", func() {

	It("bar color", func() {
		tests := []struct {
			value       int
			expectedBar string
		}{
			{value: 0, expectedBar: "[green::]▉[white::]"},
			{value: 16, expectedBar: "[orange::]▉[white::]"},
			{value: 18, expectedBar: "[red::]▉[white::]"},
		}

		for _, tt := range tests {
			Expect(getBarColor(tt.value)).To(Equal(tt.expectedBar))
		}
	})

	It("progress usage string", func() {
		tests := []struct {
			percntage          float64
			epectedUsageString string
		}{
			{percntage: 0.0, epectedUsageString: "▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉  0.00%"},
			{percntage: 10.0, epectedUsageString: "[green::]▉[white::][green::]▉[white::]▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉▉ 10.00%"},
		}

		for _, tt := range tests {
			Expect(ProgressUsageString(tt.percntage)).To(Equal(tt.epectedUsageString))
		}
	})

})

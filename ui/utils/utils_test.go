package utils

import (
	"github.com/containers/podman-tui/ui/style"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils", func() {

	It("get color name", func() {
		validColor01 := tcell.ColorRed
		validC0lor01Wants := "#ff0000"
		validColor02 := tcell.ColorBlue
		validColor02Wants := "#ff"
		invalidColor03 := tcell.Color100
		invalidcolor03Wants := "#878700"
		Expect(style.GetColorHex(validColor01)).To(Equal(validC0lor01Wants))
		Expect(style.GetColorHex(validColor02)).To(Equal(validColor02Wants))
		Expect(style.GetColorHex(invalidColor03)).To(Equal(invalidcolor03Wants))
	})

	It("empty box space", func() {
		boxColor := tcell.ColorBlue
		emptyBox := EmptyBoxSpace(boxColor)
		Expect(emptyBox.GetBackgroundColor()).To(Equal(boxColor))
		Expect(emptyBox.GetTitle()).To(Equal(""))
	})

	It("validate file name", func() {
		validFilename := "filename01"
		invalidFilename := "filename:01"
		tests := []struct {
			filename string
			wantErr  bool
		}{
			{filename: "/path/to/goog/file", wantErr: false},
			{filename: "/", wantErr: false},
			{filename: "/path/to/bad:/file", wantErr: true},
			{filename: "/path/to/bad/:file", wantErr: true},
			{filename: "/:", wantErr: true},
			{filename: ":/", wantErr: true},
		}
		for _, tt := range tests {
			if tt.wantErr {
				Expect(ValidateFileName(invalidFilename)).NotTo(BeNil())
				continue
			}
			Expect(ValidateFileName(validFilename)).To(BeNil())
		}
	})

	It("valid url", func() {
		tests := []struct {
			url   string
			valid bool
		}{
			{url: "ssh://user01@host02:22", valid: true},
			{url: "ssh://user01@host02:22/unix/1000/", valid: true},
			{url: "unix:/run/1000/podman/", valid: true},
			{url: "", valid: false},
		}
		for _, tt := range tests {
			if tt.valid {
				Expect(ValidURL(tt.url)).To(BeNil())
				continue
			}
			Expect(ValidURL(tt.url)).NotTo(BeNil())
		}
	})

	It("align string width list", func() {
		tests := []struct {
			stList        []string
			alignedStList []string
			maxWidth      int
		}{
			{stList: []string{"a", "aa", "aaa"}, alignedStList: []string{"a  ", "aa ", "aaa"}, maxWidth: 3},
			{stList: []string{"bbb", "b", "bbbb"}, alignedStList: []string{"bbb ", "b   ", "bbbb"}, maxWidth: 4},
			{stList: []string{"ccc", "c", "cc"}, alignedStList: []string{"ccc", "c  ", "cc "}, maxWidth: 3},
		}

		for _, tt := range tests {
			result, maxWidth := AlignStringListWidth(tt.stList)
			Expect(maxWidth).To(Equal(tt.maxWidth))
			for i := 0; i < len(result); i++ {
				Expect(result[i]).To(Equal(tt.alignedStList[i]))
			}
		}
	})
})

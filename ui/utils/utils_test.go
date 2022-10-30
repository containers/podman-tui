package utils

import (
	"os/user"
	"path"

	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils", func() {

	It("get color name", func() {
		validColor01 := tcell.ColorRed
		validC0lor01Wants := "red"
		validColor02 := tcell.ColorBlue
		validColor02Wants := "blue"
		invalidColor03 := tcell.Color100
		invalidcolor03Wants := ""
		Expect(GetColorName(validColor01)).To(Equal(validC0lor01Wants))
		Expect(GetColorName(validColor02)).To(Equal(validColor02Wants))
		Expect(GetColorName(invalidColor03)).To(Equal(invalidcolor03Wants))
	})

	It("empty box space", func() {
		boxColor := tcell.ColorBlue
		emptyBox := EmptyBoxSpace(boxColor)
		Expect(emptyBox.GetBackgroundColor()).To(Equal(boxColor))
		Expect(emptyBox.GetTitle()).To(Equal(""))
	})

	It("resolve home directory", func() {

		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		path01 := path.Join(user.HomeDir, "path01")
		path01Wants := path01
		path02 := "~/path02"
		path02Wants := path.Join(user.HomeDir, "path02")

		path01Resolv, err := ResolveHomeDir(path01)
		Expect(path01Resolv).To(Equal(path01Wants))
		Expect(err).To(BeNil())

		path02Resolv, err := ResolveHomeDir(path02)
		Expect(path02Resolv).To(Equal(path02Wants))
		Expect(err).To(BeNil())
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

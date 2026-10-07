package style_test

import (
	"github.com/containers/podman-tui/ui/style"
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("theme", Ordered, func() {
	It("rejects an unknown theme", func() {
		Expect(style.SetTheme("unknown")).To(MatchError(style.ErrUnknownTheme))
	})

	It("keeps the default colors", func() {
		Expect(style.SetTheme(style.ThemeDefault)).To(Succeed())
		Expect(style.GetColorHex(style.DialogFgColor)).To(Equal("#fffaf0"))
		Expect(style.GetColorName(style.PrgBarWarnColor)).To(Equal("orange"))
		Expect(style.ButtonFgColor).To(Equal(tcell.ColorWhite))
		Expect(style.TableSelectedStyle).To(Equal(tcell.StyleDefault))
	})

	It("uses the terminal's default and ANSI colors", func() {
		Expect(style.SetTheme(style.ThemeTerminal)).To(Succeed())
		Expect(style.FgColor).To(Equal(tcell.ColorDefault))
		Expect(style.BgColor).To(Equal(tcell.ColorDefault))
		Expect(style.GetColorHex(style.DialogFgColor)).To(Equal("-"))
		Expect(style.GetColorHex(style.DialogBorderColor)).To(Equal("purple"))
		Expect(style.GetColorHex(style.ButtonFgColor)).To(Equal("black"))
		Expect(style.GetColorHex(style.PageHeaderFgColor)).To(Equal("black"))
		_, _, selectedAttrs := style.TableSelectedStyle.Decompose()
		Expect(selectedAttrs & tcell.AttrReverse).To(Equal(tcell.AttrReverse))
		Expect(style.DialogSelectedStyle).To(Equal(style.TableSelectedStyle))
		Expect(style.GetColorName(style.PrgBarWarnColor)).To(Equal("olive"))
		Expect(style.GetColorName(tcell.ColorGray)).To(Equal("gray"))
		Expect(style.GetColorHex(tcell.NewRGBColor(255, 0, 0))).To(Equal("#ff0000"))
	})
})

package utils

import (
	"github.com/gdamore/tcell/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

var _ = Describe("keys", func() {

	It("key information", func() {
		nextScreenKey := uiKeyInfo{
			Key:      tcell.Key(256),
			KeyRune:  rune('l'),
			KeyLabel: "l",
			KeyDesc:  "switch to next screen",
		}
		previousScreenKey := uiKeyInfo{
			Key:      tcell.Key(256),
			KeyRune:  rune('h'),
			KeyLabel: "h",
			KeyDesc:  "switch to previous screen",
		}
		Expect(nextScreenKey.Label()).To(Equal(NextScreenKey.KeyLabel))
		Expect(nextScreenKey.Rune()).To(Equal(NextScreenKey.KeyRune))
		Expect(nextScreenKey.EventKey()).To(Equal(NextScreenKey.Key))
		Expect(nextScreenKey.Description()).To(Equal(NextScreenKey.KeyDesc))

		Expect(previousScreenKey.Label()).To(Equal(PreviousScreenKey.KeyLabel))
		Expect(previousScreenKey.Rune()).To(Equal(PreviousScreenKey.KeyRune))
		Expect(previousScreenKey.EventKey()).To(Equal(PreviousScreenKey.Key))
		Expect(previousScreenKey.Description()).To(Equal(PreviousScreenKey.KeyDesc))
	})

	It("parse key events", func() {
		zerolog.SetGlobalLevel(zerolog.Disabled)

		upEvent := tcell.NewEventKey(tcell.KeyUp, rune('k'), tcell.ModNone)
		parsedKey := ParseKeyEventKey(upEvent)
		Expect(parsedKey.Key()).To(Equal(MoveUpKey.Key))

		downEvent := tcell.NewEventKey(tcell.KeyDown, rune('j'), tcell.ModNone)
		parsedKey = ParseKeyEventKey(downEvent)
		Expect(parsedKey.Key()).To(Equal(MoveDownKey.Key))

		leftEvent := tcell.NewEventKey(tcell.KeyLeft, rune('h'), tcell.ModNone)
		parsedKey = ParseKeyEventKey(leftEvent)
		Expect(parsedKey.Key()).To(Equal(PreviousScreenKey.Key))

		rightEvent := tcell.NewEventKey(tcell.KeyRight, rune('l'), tcell.ModNone)
		parsedKey = ParseKeyEventKey(rightEvent)
		Expect(parsedKey.Key()).To(Equal(NextScreenKey.Key))

		enterEvent := tcell.NewEventKey(tcell.KeyEnter, rune(' '), tcell.ModNone)
		parsedKey = ParseKeyEventKey(enterEvent)
		Expect(parsedKey.Key()).To(Equal(tcell.KeyEnter))
	})
})

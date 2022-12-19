# tvxwidgets


[![PkgGoDev](https://pkg.go.dev/badge/github.com/navidys/tvxwidgets)](https://pkg.go.dev/github.com/navidys/tvxwidgets)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-brightgreen.svg)](https://goreportcard.com/report/github.com/navidys/tvxwidgets)

tvxwidgets provides extra widgets for [tview](https://github.com/rivo/tview).  

![Screenshot](demo.gif)

## Widgets

* [bar chart](./demos/barchart/)
* [activity mode gauge](./demos/gauge_am/)
* [percentage mode gauge](./demos/gauge_pm/)
* [utilisation mode gauge](./demos/gauge_um/)
* [message dialog (info and error)](./demos/dialog/)
* [spinner](./demos/spinner/)
* [plot (linechart, scatter)](./demos/plot/)
* [sparkline](./demos/sparkline/)


## Example

```go
package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	gauge := tvxwidgets.NewActivityModeGauge()
	gauge.SetTitle("activity mode gauge")
	gauge.SetPgBgColor(tcell.ColorOrange)
	gauge.SetRect(10, 4, 50, 3)
	gauge.SetBorder(true)

	update := func() {
		tick := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-tick.C:
				gauge.Pulse()
				app.Draw()
			}
		}
	}
	go update()

	if err := app.SetRoot(gauge, false).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

```

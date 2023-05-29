package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Page Management Sim")
	rs := MakeGRS(4, 32, PAGE_SIZE)
	gui := guiGlobalReqState(&rs)
	w.SetContent(gui)
	w.ShowAndRun()
}

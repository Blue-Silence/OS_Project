package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var i int = 1

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Container")

	s := MakeESet(20, 5)

	content2 := container.New(layout.NewHBoxLayout(), guiESet(&s))
	myWindow.SetContent(content2)
	myWindow.ShowAndRun()
}

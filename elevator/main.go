package elevator

import (
	"image/color"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var i int = 1

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Container")
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}
	red := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	blue := color.NRGBA{R: 0, B: 255, A: 255}
	cOLOR := []color.NRGBA{green, red, blue}

	text1 := canvas.NewText("Hello", green)
	//text2 := widget.NewLabel("text")
	text3 := canvas.NewText("Hello!!!", green)
	//text2.Move(fyne.NewPos(20, 20))
	//content := container.NewWithoutLayout(text1, text2)
	content := container.New(layout.NewGridLayout(3), text1, text3)
	//content11 := container.New(layout.NewGridLayout(3), text2, text2)

	content3 := widget.NewButton("click me", func() {
		log.Println("tapped", i, cOLOR[i%2])
		text1.Color = cOLOR[i%3]
		text1.Refresh()
		i++
	})

	conL1a := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText("L1", green))
	conL2a := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText("L2", green))
	conL3a := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText("L3", green))
	conL1b := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText("L1", green))
	conL2b := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText("L2", green))
	conL3b := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText("L3", green))
	conLta := container.New(layout.NewVBoxLayout(), conL1a, conL2a, conL3a)
	conLtb := container.New(layout.NewVBoxLayout(), conL1b, conL2b, conL3b)
	content2 := container.New(layout.NewHBoxLayout(), conLta, conLtb, content, content3)
	myWindow.SetContent(content2)
	myWindow.ShowAndRun()
}

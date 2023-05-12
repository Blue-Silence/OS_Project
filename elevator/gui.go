package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var green color.NRGBA = color.NRGBA{R: 0, G: 180, B: 0, A: 255}
var gray color.NRGBA = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
var orange color.NRGBA = color.NRGBA{R: 225, G: 128, B: 0, A: 255}

func guiESet(s *ESet) *fyne.Container {
	go func() {
		for {
			sec := 1
			time.Sleep(time.Duration(sec) * time.Second)
			for _, v := range s.clockChs {
				v <- 0
			}
		}
	}()
	p := guiPannel(s)
	var objs []fyne.CanvasObject
	objs = append(objs, p)
	for i, _ := range s.es {
		objs = append(objs, guiE(&s.es[i]))
	}
	return container.New(layout.NewHBoxLayout(), objs...) //, container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText("L3", green)))
}

func guiE(e *ECB) *fyne.Container {
	//currentFloor := canvas.NewText("   1  ", green)
	//currentState := canvas.NewText("   STAY  ", green)
	currentFloor := widget.NewLabel("   1     ")
	currentState := widget.NewLabel("   STAY  ")
	currentFloor.TextStyle = fyne.TextStyle{Monospace: true}
	currentState.TextStyle = fyne.TextStyle{Monospace: true}
	var buttons []fyne.CanvasObject
	var bs []*widget.Button
	var bcs []*canvas.Rectangle
	for i := e.topFloor; i > 0; i-- {
		i2 := i - 1
		b := widget.NewButton(fmt.Sprintf("%v", i), func() {
			//p.setTarget(Upward, i2, true)
			e.insertInternalTarget(i2)
		})
		b.Importance = widget.LowImportance
		bc := canvas.NewRectangle(color.White)

		bs = append(bs, b)
		bcs = append(bcs, bc)
		buttons = append(buttons, container.New(layout.NewMaxLayout(), bc, b))

	}

	go func() {
		for {
			e.mu.Lock()
			switch {
			case e.State == Run && e.Dir == Upward:
				currentState.SetText("  Going up  ")
			case e.State == Run && e.Dir == Downward:
				currentState.SetText(" Going down ")
			case e.State == Stay1:
				currentState.SetText("Door opening")
			case e.State == Stay2:
				currentState.SetText("  Door open ")
			case e.State == Stay3:
				currentState.SetText("Door closing")
			case e.State == Idle:
				currentState.SetText("    STAY    ")
			default:
				currentState.SetText("What???")
			}
			currentFloor.SetText(fmt.Sprintf("%3d", e.floor+1))
			for i := 0; i < e.topFloor; i++ {
				if e.internalButton[i] {
					//r, _ := fyne.LoadResourceFromPath("./Resource/2.png")
					//bs[e.topFloor-1-i].SetText(fmt.Sprintf("Comming:%v!", i+1))
					bcs[e.topFloor-1-i].FillColor = (orange)
				} else {
					bs[e.topFloor-1-i].SetText(fmt.Sprintf("%v", i+1))
					bcs[e.topFloor-1-i].FillColor = (gray)
				}
				bcs[e.topFloor-1-i].Refresh()
			}
			e.mu.Unlock()
			_ = <-e.signalCh
		}

	}()

	go func() {
		log.Println("Ok.....")
		for {
			_ = <-e.clockCh
			e.stateForward()
		}
	}()

	Screen := container.New(layout.NewHBoxLayout(), currentFloor, currentState)
	//Screen := container.New(layout.NewVBoxLayout(), currentFloor, currentState)

	open := widget.NewButton("", func() {
		e.doorOpen()
	})
	open.Importance = widget.LowImportance
	ro, _ := fyne.LoadResourceFromPath("./Resource/OPEN.png")
	open.SetIcon(ro)
	close := widget.NewButton("", func() {
	})
	close.Importance = widget.LowImportance
	rc, _ := fyne.LoadResourceFromPath("./Resource/CLOSE.png")
	close.SetIcon(rc)
	//doorP := container.New(layout.NewHBoxLayout(), open, close)
	doorP := container.New(layout.NewGridLayout(2), open, close)

	buttons = append(buttons, Screen, doorP)
	//buttons = append(buttons, doorP)
	//buttons = append([]fyne.CanvasObject{Screen}, buttons...)

	//buttons = append(buttons, Screen, open, close)

	return container.New(layout.NewVBoxLayout(), buttons...)
}

func guiPannel(s *ESet) *fyne.Container {
	var buttons []fyne.CanvasObject
	var ups []*canvas.Rectangle
	var downs []*canvas.Rectangle
	p := s.pannel
	for i := p.topFloor; i > 0; i-- {
		i2 := i - 1
		up := widget.NewButton("", func() {
			//p.setTarget(Upward, i2, true)
			s.requestElevator(Upward, i2)
		})
		up.Importance = widget.LowImportance
		ru, _ := fyne.LoadResourceFromPath("./Resource/UP.png")
		up.SetIcon(ru)
		up.IconPlacement = widget.ButtonIconLeadingText
		upR := canvas.NewRectangle(color.White)
		ups = append(ups, upR)

		down := widget.NewButton("", func() {
			s.requestElevator(Downward, i2)
		})
		down.Importance = widget.LowImportance
		rd, _ := fyne.LoadResourceFromPath("./Resource/DOWN.png")
		down.SetIcon(rd)
		down.IconPlacement = widget.ButtonIconLeadingText
		downR := canvas.NewRectangle(color.White)
		downs = append(downs, downR)
		//tag := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText(fmt.Sprintf("%3v", i), green))
		tag := canvas.NewText(fmt.Sprintf("%3v", i), color.Black)
		tag.TextStyle = fyne.TextStyle{Monospace: true}
		tag.TextSize = 20.0
		buttons = append(buttons, container.New(layout.NewHBoxLayout(), tag, container.New(layout.NewMaxLayout(), upR, up), container.New(layout.NewMaxLayout(), downR, down)))
	}

	go func() {
		for {
			_ = <-p.signalCh
			p.mu.Lock()
			//log.Println("------------------------------------------------------------")
			//log.Println(p.upTarget)
			//log.Println(p.downTarget)
			for i := 0; i < p.topFloor; i++ {
				if p.upTarget[i] {
					//r, _ := fyne.LoadResourceFromPath("./Resource/ButtonUp2.png")
					ups[p.topFloor-1-i].FillColor = (green)
					//ups[p.topFloor-1-i].SetText("Upping!")
				} else {
					//r, _ := fyne.LoadResourceFromPath("./Resource/UP.png")
					ups[p.topFloor-1-i].FillColor = (color.White)
					//ups[p.topFloor-1-i].SetText("Up    !")
				}
				if p.downTarget[i] {
					downs[p.topFloor-1-i].FillColor = (green)
				} else {
					//r, _ := fyne.LoadResourceFromPath("./Resource/DOWN.png")
					downs[p.topFloor-1-i].FillColor = (color.White)
				}
				ups[p.topFloor-1-i].Refresh()
				downs[p.topFloor-1-i].Refresh()
			}
			p.mu.Unlock()
			//log.Println("############################################################")
			//log.Println(p.upTarget)
			//log.Println(p.downTarget)
			//log.Println("------------------------------------------------------------")

		}
	}()

	return container.New(layout.NewVBoxLayout(), buttons...)
}

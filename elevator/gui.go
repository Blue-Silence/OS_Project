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

func guiESet(s *ESet) *fyne.Container {
	go func() {
		//log.Println("HELLO!!!")
		//secCount := 0
		for {
			sec := 1
			time.Sleep(time.Duration(sec) * time.Second)
			//log.Println(secCount)
			//secCount++
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
	currentFloor := widget.NewLabel("   1  ")
	currentState := widget.NewLabel("   STAY  ")
	var buttons []fyne.CanvasObject
	var bs []*widget.Button
	for i := e.topFloor; i > 0; i-- {
		i2 := i - 1
		b := widget.NewButton(fmt.Sprintf(" %3v ", i), func() {
			//p.setTarget(Upward, i2, true)
			e.insertInternalTarget(i2)
		})
		bs = append(bs, b)
		buttons = append(buttons, b)
	}

	go func() {
		for {
			e.mu.Lock()
			switch {
			case e.State == Run && e.Dir == Upward:
				currentState.SetText("Going up")
			case e.State == Run && e.Dir == Downward:
				currentState.SetText("Going down")
			case e.State == Stay1:
				currentState.SetText("Door opening")
			case e.State == Stay2:
				currentState.SetText("Door open")
			case e.State == Stay3:
				currentState.SetText("Door closing")
			case e.State == Idle:
				currentState.SetText("NOTHING")
			default:
				currentState.SetText("What???")
			}
			currentFloor.SetText(fmt.Sprintf("%3d", e.floor+1))
			for i := 0; i < e.topFloor; i++ {
				if e.internalButton[i] {
					bs[e.topFloor-1-i].SetText(fmt.Sprintf("Comming:%v!", i+1))
				} else {
					bs[e.topFloor-1-i].SetText(fmt.Sprintf("%v", i+1))
				}
			}
			e.mu.Unlock()
			_ = <-e.signalCh
		}

	}()

	go func() {
		log.Println("Ok.....")
		for {
			//log.Println("Before!")
			_ = <-e.clockCh
			//log.Println("State changing!")
			e.stateForward()
		}
	}()

	//Screen := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), container.New(layout.NewHBoxLayout(), currentFloor, currentState))
	Screen := container.New(layout.NewHBoxLayout(), currentFloor, currentState)
	buttons = append(buttons, Screen)

	return container.New(layout.NewVBoxLayout(), buttons...)
}

func guiPannel(s *ESet) *fyne.Container {
	var buttons []fyne.CanvasObject
	var ups []*widget.Button
	var downs []*widget.Button
	p := s.pannel
	for i := p.topFloor; i > 0; i-- {
		i2 := i - 1
		up := widget.NewButton("Up    !", func() {
			//p.setTarget(Upward, i2, true)
			s.requestElevator(Upward, i2)
		})
		ups = append(ups, up)
		down := widget.NewButton("Down   !", func() {
			//p.setTarget(Downward, i2, true)
			s.requestElevator(Downward, i2)
		})
		downs = append(downs, down)
		tag := container.New(layout.NewMaxLayout(), canvas.NewRectangle(color.White), canvas.NewText(fmt.Sprintf(" %3v ", i), green))
		buttons = append(buttons, container.New(layout.NewHBoxLayout(), tag, up, down))
	}

	go func() {
		for {
			_ = <-p.signalCh
			p.mu.Lock()
			log.Println("------------------------------------------------------------")
			log.Println(p.upTarget)
			log.Println(p.downTarget)
			for i := 0; i < p.topFloor; i++ {
				if p.upTarget[i] {
					ups[p.topFloor-1-i].SetText("Upping!")
				} else {
					ups[p.topFloor-1-i].SetText("Up    !")
				}
				if p.downTarget[i] {
					downs[p.topFloor-1-i].SetText("Downing!")
				} else {
					downs[p.topFloor-1-i].SetText("Down   !")
				}
			}
			p.mu.Unlock()
			log.Println("############################################################")
			log.Println(p.upTarget)
			log.Println(p.downTarget)
			log.Println("------------------------------------------------------------")

		}
	}()

	return container.New(layout.NewVBoxLayout(), buttons...)
}

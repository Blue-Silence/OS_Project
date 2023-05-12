package main

import (
	"fmt"
	"image/color"
	"log"
	"sync"
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

type GlobalReqState struct {
	currentPolicy int
	timeInterval  int
	msgCh         chan string
	mu            sync.Mutex
}

func (g *GlobalReqState) setFollowingPolicy(policy int) {
	g.mu.Lock()
	g.currentPolicy = policy
	g.mu.Unlock()
	g.msgCh <- fmt.Sprintln("Set policy to:", policy)
}

func (g *GlobalReqState) guiGlobalState(s *GlobalState) *fyne.Container {
	var lastEnters []*widget.Label
	var lastUseds []*widget.Label
	var vailds []*widget.Label
	var vns []*widget.Label
	var rects []*canvas.Rectangle

	var pps []fyne.CanvasObject
	for pn, _ := range s.physicalPs {
		lastEnter := widget.NewLabel("0")
		lastUsed := widget.NewLabel("0")
		vn := widget.NewLabel("0")
		vaild := widget.NewLabel("False")
		rect := canvas.NewRectangle(color.White)

		lastEnters = append(lastEnters, lastEnter)
		lastUseds = append(lastUseds, lastUsed)
		vns = append(vns, vn)
		vailds = append(vailds, vaild)
		rects = append(rects, rect)

		ppT := container.New(layout.NewVBoxLayout(),
			widget.NewLabel(fmt.Sprint("Physical page number:", pn)),
			container.New(layout.NewHBoxLayout(), widget.NewLabel("Virtual page number:"), vn),
			container.New(layout.NewHBoxLayout(), widget.NewLabel("Cycle since last used:"), lastUsed),
			container.New(layout.NewHBoxLayout(), widget.NewLabel("Cycle since page in:"), lastEnter),
			container.New(layout.NewHBoxLayout(), widget.NewLabel("Vaild:"), vaild))
		pp := container.New(layout.NewMaxLayout(), rect, ppT)
		pps = append(pps, pp)
	}

	ins := canvas.NewText("\n\n\n", orange)

	go func() {
		for {
			msg := <-s.signalCh
			s.mu.Lock()
			switch {
			case msg.isReset:
				ins.Text = fmt.Sprintf("%v%v\n", ins.Text, "Reset")
				ins.Refresh()
				for a, _ := range s.physicalPs {
					s.physicalPs[a] = nil
					lastEnters[msg.PN].SetText("0")
					lastUseds[msg.PN].SetText("0")
					vailds[msg.PN].SetText("False")
					vns[msg.PN].SetText("0")
					rects[msg.PN].FillColor = color.White
					rects[msg.PN].Refresh()
				}
			case msg.isHit:
				ins.Text = fmt.Sprintf("%v%v\n", ins.Text, fmt.Sprintf("Page hit for request at:%v", msg.reqAddress))
				ins.Refresh()
				t := g.timeInterval
				if t > 0 {
					rects[msg.PN].FillColor = green
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t * 100))
					rects[msg.PN].FillColor = color.White
					rects[msg.PN].Refresh()
				}
			case msg.isReplace:
				ins.Text = fmt.Sprintf("%v%v\n", ins.Text, fmt.Sprintf("Page miss for request at:%v, replacing in at:%v", msg.reqAddress, msg.PN))
				ins.Refresh()
				t := g.timeInterval
				if t > 0 {
					rects[msg.PN].FillColor = orange
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t * 100))
					rects[msg.PN].FillColor = color.White
					rects[msg.PN].Refresh()
				}

				lastEnters[msg.PN].SetText("0")
				lastUseds[msg.PN].SetText("0")
				vns[msg.PN].SetText(fmt.Sprint(msg.VN))
				if t > 0 {
					rects[msg.PN].FillColor = green
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t * 100))
					rects[msg.PN].FillColor = color.White
					rects[msg.PN].Refresh()
				}
			case !msg.isReplace && msg.isHit:
				ins.Text = fmt.Sprintf("%v%v\n", ins.Text, fmt.Sprintf("Page miss for request at:%v, switching in at:%v", msg.reqAddress, msg.PN))
				ins.Refresh()
				t := g.timeInterval
				if t > 0 {
					rects[msg.PN].FillColor = orange
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t * 100))
					rects[msg.PN].FillColor = color.White
					rects[msg.PN].Refresh()
				}

				lastEnters[msg.PN].SetText("0")
				lastUseds[msg.PN].SetText("0")
				vailds[msg.PN].SetText("True")
				vns[msg.PN].SetText(fmt.Sprint(msg.VN))
				if t > 0 {
					rects[msg.PN].FillColor = green
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t * 100))
					rects[msg.PN].FillColor = color.White
					rects[msg.PN].Refresh()
				}
			default:
				log.Fatal("!!!!!")
			}
		}
	}()

	return container.New(layout.NewVBoxLayout(), append(pps, container.NewVScroll(ins))...)
}

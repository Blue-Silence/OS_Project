package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	SEQ    int = 0
	RANDOM int = 1
)

var green color.NRGBA = color.NRGBA{R: 0, G: 180, B: 0, A: 255}
var gray color.NRGBA = color.NRGBA{R: 160, G: 160, B: 160, A: 255}
var orange color.NRGBA = color.NRGBA{R: 225, G: 128, B: 0, A: 255}
var sec = time.Duration(1000) * time.Millisecond

type GlobalReqState struct {
	policy       int
	timeInterval int
	lastAddr     int

	s GlobalState

	strategy   int
	upperBound int

	ppageN, vpageN int

	msgCh chan string
	mu    sync.Mutex
}

func (G *GlobalReqState) reset() {
	G.mu.Lock()
	G.policy = FIFO
	G.timeInterval = 0
	G.lastAddr = -1
	G.strategy = SEQ
	oldCh := G.msgCh
	G.msgCh = make(chan string, 1024)
	G.mu.Unlock()
	oldCh <- "Reset ok!"
}

func MakeGRS(ppageN int, vpageN int, pageSize int) GlobalReqState {
	var G GlobalReqState
	G.msgCh = make(chan string, 1024)
	G.reset()
	G.s.signalCh = make(chan ReqMsg, 2)
	G.s.reset(ppageN, vpageN, FIFO)
	G.ppageN = ppageN
	G.vpageN = vpageN
	G.upperBound = vpageN*pageSize - 1
	return G
}

func (g *GlobalReqState) setFollowingPolicy(policy int) {
	g.mu.Lock()
	g.policy = policy
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
		rect := canvas.NewRectangle(gray)

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

	ins := widget.NewLabel("                                                   \n")
	ins.TextStyle = fyne.TextStyle{Monospace: true}
	//ins2 := widget.NewTextGridFromString("                                                   \n")
	//ins2.TextStyle = fyne.TextStyle{Monospace: true}

	insR := canvas.NewRectangle(gray)
	//ins.SetMinSize(fyne.Size{4, 400})
	reqC := canvas.NewText("0", orange)
	missC := canvas.NewText("0", orange)
	outputWindow := container.NewVScroll(container.New(layout.NewMaxLayout(), insR, ins))
	outputWindow.SetMinSize(fyne.Size{4, 400})
	globalStatus := container.New(layout.NewVBoxLayout(), container.New(layout.NewHBoxLayout(), canvas.NewText("Request", orange), reqC), container.New(layout.NewHBoxLayout(), canvas.NewText("Miss:", orange), missC), outputWindow)
	go func() {
		for {
			s.mu.Lock()
			signalCh := s.signalCh
			//msg := <-s.signalCh
			s.mu.Unlock()
			msg := <-signalCh
			s.mu.Lock()
			insText := ""
			switch {

			case msg.isReset:
				insText = fmt.Sprintf("%v\n", "Reset                                              ")
				ins.Refresh()
				//ins2.SetText("Reset                                                                ")
				for a, _ := range s.physicalPs {
					s.physicalPs[a] = nil
					lastEnters[a].SetText("0")
					lastUseds[a].SetText("0")
					vailds[a].SetText("False")
					vns[a].SetText("0")
					rects[a].FillColor = gray //color.White
					rects[a].Refresh()
				}
			case msg.isHit:
				insText = fmt.Sprintf("%v%v\n", ins.Text, fmt.Sprintf("Page hit for request at:%v", msg.reqAddress))
				ins.Refresh()
				//ins2.SetText(fmt.Sprintf("%v%v\n", ins2.Text(), fmt.Sprintf("Page hit for request at:%v", msg.reqAddress)))
				t := g.timeInterval
				if t > 0 {
					rects[msg.PN].FillColor = green
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t*100) * time.Millisecond)
					rects[msg.PN].FillColor = gray //color.White
					rects[msg.PN].Refresh()
				}
			case msg.isReplace:
				insText = fmt.Sprintf("%v%v\n", ins.Text, fmt.Sprintf("Page miss for request at:%v, replacing in at:%v", msg.reqAddress, msg.PN))
				ins.Refresh()
				//ins2.SetText(fmt.Sprintf("%v%v\n", ins2.Text(), fmt.Sprintf("Page miss for request at:%v, replacing in at:%v", msg.reqAddress, msg.PN)))
				t := g.timeInterval
				if t > 0 {
					rects[msg.PN].FillColor = orange
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t*100) * time.Millisecond)
					rects[msg.PN].FillColor = gray //color.White
					rects[msg.PN].Refresh()
				}

				lastEnters[msg.PN].SetText("0")
				lastUseds[msg.PN].SetText("0")
				vns[msg.PN].SetText(fmt.Sprint(msg.VN))
				if t > 0 {
					rects[msg.PN].FillColor = green
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t*100) * time.Millisecond)
					rects[msg.PN].FillColor = gray //color.White
					rects[msg.PN].Refresh()
				}
			case !msg.isReplace && !msg.isHit:
				insText = fmt.Sprintf("%v%v\n", ins.Text, fmt.Sprintf("Page miss for request at:%v, switching in at:%v", msg.reqAddress, msg.PN))

				//ins2.SetText(fmt.Sprintf("%v%v\n", ins2.Text(), fmt.Sprintf("Page miss for request at:%v, switching in at:%v", msg.reqAddress, msg.PN)))
				t := g.timeInterval
				if t > 0 {
					rects[msg.PN].FillColor = orange
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t*100) * time.Millisecond)
					rects[msg.PN].FillColor = gray //color.White
					rects[msg.PN].Refresh()
				}

				lastEnters[msg.PN].SetText("0")
				lastUseds[msg.PN].SetText("0")
				vailds[msg.PN].SetText("True")
				vns[msg.PN].SetText(fmt.Sprint(msg.VN))
				if t > 0 {
					rects[msg.PN].FillColor = green
					rects[msg.PN].Refresh()
					time.Sleep(time.Duration(t*100) * time.Millisecond)
					rects[msg.PN].FillColor = gray //color.White
					rects[msg.PN].Refresh()
				}
			default:
				log.Fatal("!!!!!")
			}

			for a, _ := range msg.physicalPs {
				if s.physicalPs[a] != nil {

					lastEnters[a].SetText(fmt.Sprint(s.physicalPs[a].lastEnter))
					lastUseds[a].SetText(fmt.Sprint(s.physicalPs[a].lastUsed))
					//vailds[a].SetText("False")
					//vns[a].SetText("0")
					//rects[a].FillColor = gray //color.White
					//rects[a].Refresh()
				}
			}

			split := strings.SplitAfter(insText, "\n")
			if len(split) > 30 {
				insText = ""
				for _, s := range split[1:] {
					insText = fmt.Sprintf("%v%v", insText, s)
				}
			}

			ins.Text = insText
			ins.Refresh()

			missC.Text = (fmt.Sprint(msg.missCounter))
			reqC.Text = (fmt.Sprint(msg.reqCounter))
			missC.Refresh()
			reqC.Refresh()
			outputWindow.ScrollToBottom()
			s.mu.Unlock()
		}
	}()

	return container.New(layout.NewHBoxLayout(), append(pps, globalStatus)...)
}

func guiGlobalReqState(g *GlobalReqState) *fyne.Container {

	a1 := func() {
		g.mu.Lock()
		switch g.strategy {
		case SEQ:
			g.lastAddr = (g.lastAddr + 1) % (g.upperBound + 1)
		case RANDOM:
			g.lastAddr = rand.Int() % (g.upperBound + 1)
		}

		g.s.reqAddress(g.lastAddr, g.policy)
		g.mu.Unlock()
	}

	addNote := widget.NewLabel("   Adding instruction   ")
	addNote.TextStyle = fyne.TextStyle{Monospace: true}
	add1 := widget.NewButton("Add 1 instr", func() {
		a1()
		g.msgCh <- fmt.Sprintf("Adding access to %v", g.lastAddr)
	})

	add50 := widget.NewButton("Add 50 instr", func() {
		for i := 0; i < 50; i++ {
			a1()
		}
		//a1()
		g.msgCh <- fmt.Sprintf("Adding 50 access", g.lastAddr)
	})

	inputInstrNum := widget.NewEntry()
	inputInstrNum.SetPlaceHolder("Enter instruction number.")
	addN := widget.NewButton("      Add      ", func() {
		n := 0
		fmt.Sscanf(inputInstrNum.Text, "%v", &n)
		for i := 0; i < n; i++ {
			a1()
		}
		g.msgCh <- fmt.Sprintf("Adding %v access", n)
	})

	instrInsertWin := container.New(layout.NewVBoxLayout(), addNote, add1, add50, container.New(layout.NewVBoxLayout(), inputInstrNum, addN))

	addrNote := widget.NewLabel("      Set insert address      ")
	addrNote.TextStyle = fyne.TextStyle{Monospace: true}
	setSEQ := widget.NewButton("SET SEQ", func() {
		g.mu.Lock()
		g.strategy = SEQ
		g.mu.Unlock()
		g.msgCh <- fmt.Sprintf("Change stratgy to SEQ.")
	})

	setRAN := widget.NewButton("          SET RANDOM          ", func() {
		g.mu.Lock()
		g.strategy = RANDOM
		g.mu.Unlock()
		g.msgCh <- fmt.Sprintf("Change stratgy to RANDOM.")
	})

	inputAddr := widget.NewEntry()
	inputAddr.SetPlaceHolder("Enter next instr addr.")
	setAddr := widget.NewButton("Set", func() {
		g.mu.Lock()
		n := 0
		fmt.Sscanf(inputAddr.Text, "%v", &n)
		if n < 0 {
			n = 0
		} else {
			n = n % (g.upperBound + 1)
		}
		g.lastAddr = n - 1
		g.mu.Unlock()
		g.msgCh <- fmt.Sprintf("Set next addr at %v", n)
	})

	container.NewWithoutLayout()
	addrWin := container.New(layout.NewVBoxLayout(), addrNote, setSEQ, setRAN, container.New(layout.NewVBoxLayout(), inputAddr, setAddr))

	policyNote := widget.NewLabel("      Set replace policy      ")
	policyNote.TextStyle = fyne.TextStyle{Monospace: true}
	setFIFO := widget.NewButton("          SET FIFO          ", func() {
		g.mu.Lock()
		g.policy = FIFO
		g.mu.Unlock()
		g.msgCh <- fmt.Sprintf("Change policy to FIFO.")
	})

	setLRU := widget.NewButton("SET LRU", func() {
		g.mu.Lock()
		g.policy = LRU
		g.mu.Unlock()
		g.msgCh <- fmt.Sprintf("Change policy to LRU.")
	})
	replacePWin := container.New(layout.NewVBoxLayout(), policyNote, setFIFO, setLRU)

	globalNote := widget.NewLabel("      Set motion speed and reset      ")
	globalNote.TextStyle = fyne.TextStyle{Monospace: true}
	reset := widget.NewButton("CLEAR", func() {
		g.reset()
		log.Println("OK?")
		g.mu.Lock()
		g.s.reset(g.ppageN, g.vpageN, FIFO)
		g.mu.Unlock()
		log.Println("OK!")
		g.msgCh <- fmt.Sprintf("          Clear          ")
	})

	inputTime := widget.NewEntry()
	inputTime.SetPlaceHolder("Enter time interval (*100ms).")
	setTime := widget.NewButton("Set", func() {
		g.mu.Lock()
		n := 0
		fmt.Sscanf(inputTime.Text, "%v", &n)
		if n < 0 {
			n = 0
		} else {
		}
		g.timeInterval = n
		g.mu.Unlock()
		g.msgCh <- fmt.Sprintf("Set time interval as %v*100ms", n)
	})
	gActWin := container.New(layout.NewVBoxLayout(), container.New(layout.NewVBoxLayout(), globalNote, inputTime, setTime), reset)

	ins := widget.NewLabel("                           \n")
	outputWindow := container.NewVScroll(ins)
	outputWindow.SetMinSize(fyne.Size{4, 400})
	poS := widget.NewLabel("FIFO")
	timeIS := widget.NewLabel("0")
	nextAS := widget.NewLabel("0")
	strategyS := widget.NewLabel("SEQ")
	status := container.New(layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(), widget.NewLabel("Policy:"), poS),
		container.New(layout.NewHBoxLayout(), widget.NewLabel("Time Interval:"), timeIS),
		container.New(layout.NewHBoxLayout(), widget.NewLabel("Next access address policy:"), strategyS),
		container.New(layout.NewHBoxLayout(), widget.NewLabel("Next access address:"), nextAS),
		outputWindow)

	go func() {
		for {
			g.mu.Lock()
			msgCh := g.msgCh
			g.mu.Unlock()
			str := <-msgCh
			g.mu.Lock()
			ins.Text = fmt.Sprintf("%v%v\n", ins.Text, str)
			ins.Refresh()
			switch g.strategy {
			case RANDOM:
				strategyS.SetText("RANDOM")
			case SEQ:
				strategyS.SetText("SEQ")
			}
			switch g.policy {
			case FIFO:
				poS.SetText("FIFO")
			case LRU:
				poS.SetText("LRU")
			}
			timeIS.SetText(fmt.Sprint(g.timeInterval, "*100ms"))
			nextAS.SetText(fmt.Sprint(g.lastAddr + 1))
			outputWindow.ScrollToBottom()
			g.mu.Unlock()
		}
	}()

	return container.New(layout.NewVBoxLayout(),
		g.guiGlobalState(&g.s),
		container.New(layout.NewHBoxLayout(), instrInsertWin, addrWin, replacePWin, gActWin, status))
}

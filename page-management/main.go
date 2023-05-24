package main

import (
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")
	log.Println("333")
	w.SetContent(widget.NewLabel("Hello World!"))
	log.Println("555")
	rs := MakeGRS(4, 32, PAGE_SIZE)
	//rs.s.reset(4, 10, FIFO)

	gui := guiGlobalReqState(&rs)
	//foo(&rs)
	w.SetContent(gui)
	/*go func() {
		i := 0
		for {
			rs.timeInterval = 10
			time.Sleep(time.Duration(2*1000) * time.Millisecond)
			log.Println("Hello~")
			//time.Sleep(time.Duration(5 * 10000))
			foo(&rs, i, true)
			i = (i + 1) % 100
			log.Println("23333")
			//time.Sleep(time.Duration(5 * 1000))
		}
	}()*/
	w.ShowAndRun()
}

func foo(g *GlobalReqState, seq int, f bool) {
	if !f {
		//g.s.reqAddress(rand.Int()%10*10 + rand.Int()%10)
	} else {
		g.s.reqAddress(seq, FIFO)
	}
	//time.Sleep(time.Duration(5 * 100000))
}

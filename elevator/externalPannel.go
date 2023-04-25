package main

import "sync"

type Pannel struct {
	upTarget   []bool
	downTarget []bool
	signalCh   chan int
	topFloor   int
	mu         sync.Mutex
}

func MakePannel(floor int) Pannel {
	var p Pannel
	p.upTarget = make([]bool, floor)
	p.downTarget = make([]bool, floor)
	p.signalCh = make(chan int)
	p.topFloor = floor
	return p
}

func (p *Pannel) setTarget(Dir int, floor int, v bool) bool {
	var r bool
	p.mu.Lock()
	r = p.upTarget[floor]
	switch Dir {
	case Upward:
		p.upTarget[floor] = v
	case Downward:
		p.downTarget[floor] = v
	}
	p.mu.Unlock()
	p.signalCh <- 0
	return r
}

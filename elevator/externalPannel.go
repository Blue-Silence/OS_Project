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

	switch Dir {
	case Upward:
		r = p.upTarget[floor]
		p.upTarget[floor] = v
	case Downward:
		r = p.downTarget[floor]
		p.downTarget[floor] = v
	}
	p.mu.Unlock()
	p.signalCh <- 0
	return r
}

func (p *Pannel) clearTarget(Dir int, floor int) {

	p.mu.Lock()

	switch {
	case Dir == Upward && p.upTarget[floor]:
		p.upTarget[floor] = false
	case Dir == Downward && p.downTarget[floor]:
		p.downTarget[floor] = false
	default:
		p.upTarget[floor] = false
		p.downTarget[floor] = false
	}
	p.mu.Unlock()
	p.signalCh <- 0
}

package elevator

import "sync"

type Pannel struct {
	upTarget   []bool
	downTarget []bool
	signalCh   chan int
	mu         sync.Mutex
}

func MakePannel(floor int) Pannel {
	var p Pannel
	p.upTarget = make([]bool, floor)
	p.downTarget = make([]bool, floor)
	p.signalCh = make(chan int)
	return p
}

func (p *Pannel) setTarget(Dir int, floor int, v bool) {
	p.mu.Lock()
	switch Dir {
	case Upward:
		p.upTarget[floor] = v
	case Downward:
		p.downTarget[floor] = v
	}
	p.mu.Unlock()
	p.signalCh <- 0
}

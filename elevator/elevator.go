package main

import (
	"log"
	"sync"
)

type ECB struct {
	State int
	Dir   int
	floor int
	//upTarget int[]
	//downTarget int[]
	Target         []int
	internalButton []bool
	mu             sync.Mutex
	pannel         *Pannel
	topFloor       int
	clockCh        chan int
	signalCh       chan int
}

func MakeECB(floors int, p *Pannel) ECB {
	var e ECB
	e.State = Idle
	e.floor = 0
	e.Dir = Upward
	e.Target = make([]int, floors)
	e.internalButton = make([]bool, floors)
	e.pannel = p
	e.topFloor = floors
	e.clockCh = make(chan int)
	e.signalCh = make(chan int)
	return e
}

func (e *ECB) insertTarget(f int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Target[f]++
}

func (e *ECB) insertInternalTarget(f int) {
	e.mu.Lock()
	//defer e.singalCh <- 0
	//defer e.mu.Unlock()
	if e.Target[f] == 0 {
		e.Target[f]++
	}
	e.internalButton[f] = true

	e.mu.Unlock()
	e.signalCh <- 0
}

func (e *ECB) stateForward() {
	e.mu.Lock()
	switch e.State {
	case Idle:
		e.stateForwardIdle()
	case Run:
		e.stateForwardRun()
	case Stay1:
		e.stateForwardStay1()
	case Stay2:
		e.stateForwardStay2()
	case Stay3:
		e.stateForwardStay3()
	}
	e.mu.Unlock()
	e.signalCh <- 0

}

func (e *ECB) distanceCal(dir int, floor int) int {
	r := 0
	e.mu.Lock()
	defer e.mu.Unlock()
	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}
	targetCount := 0
	upperBound := 0
	lowerBound := e.topFloor - 1
	for i, v := range e.Target {
		if v > 0 {
			targetCount++
			if i > upperBound {
				upperBound = i
			}
			if i < lowerBound {
				lowerBound = i
			}
		}
	}

	switch e.State {
	case Idle:
		r = abs(floor - e.floor)
	//case Run:
	default:
		switch {
		case (e.floor < floor && e.Dir == Upward) || (e.floor > floor && e.Dir == Downward):
			r = abs(floor - e.floor)
		case e.Dir == Upward:
			r = abs(upperBound-e.floor) + abs(upperBound-floor)
		case e.Dir == Downward:
			r = abs(lowerBound-e.floor) + abs(lowerBound-floor)
		}
	}

	return r + targetCount*2 // Take into account the current load

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (e *ECB) stateForwardIdle() {
	upCount, downCount := e.targetCount()

	switch {
	case e.Target[e.floor] > 0:
		e.Target[e.floor]--
		e.internalButton[e.floor] = false
		e.State = Stay1
	case upCount+downCount == 0:
	case upCount >= downCount:
		e.State = Run
		e.Dir = Upward
	case upCount < downCount:
		e.State = Run
		e.Dir = Downward
	}
	log.Println("Idle")
}

func (e *ECB) stateForwardRun() {
	switch e.Dir {
	case Upward:
		e.floor++
		switch {
		case e.Target[e.floor] > 0:
			e.State = Stay1
		case e.floor == e.topFloor-1:
			e.State = Stay3
		}
	case Downward:
		e.floor--
		switch {
		case e.Target[e.floor] > 0:
			e.State = Stay1
		case e.floor == 0:
			e.State = Stay3
		}
	}
	log.Println("e.floor:", e.floor, " e.Target:", e.Target)
}

func (e *ECB) stateForwardStay1() {
	e.stateToStay2()
}

func (e *ECB) stateForwardStay2() {
	switch {
	case e.Target[e.floor] > 0:
		e.stateToStay2()
		//e.Target[e.floor]--
	default:
		e.State = Stay3
		//e.Target[e.floor]--
	}
	//e.Target[e.floor]--
	e.internalButton[e.floor] = false
}

func (e *ECB) stateForwardStay3() {
	switch {
	case e.Target[e.floor] > 0:
		e.stateToStay2()
		//e.Target[e.floor]--
	default:
		upCount, downCount := e.targetCount()
		switch {
		case e.Dir == Upward && upCount > 0 || e.Dir == Downward && downCount > 0:
			e.State = Run
		default:
			e.State = Idle
			e.stateForwardIdle()
		}
	}
	//e.Target[e.floor]--
	e.internalButton[e.floor] = false
}

func (e *ECB) stateToStay2() {
	e.State = Stay2
	if e.Target[e.floor] > 0 {
		e.Target[e.floor]--
	}
	//e.Target[e.floor]--
	e.internalButton[e.floor] = false
	//And also do something to clear the external button. TO BE DONE.
	dir := e.Dir
	f := e.floor
	e.mu.Unlock()
	log.Println("HELLO!")
	log.Println("dir:", dir, "  f:", f)
	e.pannel.clearTarget(dir, f)
	//if !e.pannel.setTarget(dir, f, false) {
	//	log.Println("OKKKKK")
	//	e.pannel.setTarget(reverseDir(dir), f, false)
	//}
	e.mu.Lock()

}

func (e *ECB) targetCount() (int, int) {
	upCount := 0
	downCount := 0
	for i, v := range e.Target {
		switch {
		case v == 0:
		case i > e.floor:
			upCount++
		case i < e.floor:
			downCount++
		}
	}
	return upCount, downCount
}

func reverseDir(dir int) int {
	switch dir {
	case Upward:
		return Downward
	case Downward:
		return Upward
	}
	log.Fatal("Warning!")
	return Upward
}

func (e *ECB) doorOpen() {
	e.mu.Lock()
	switch e.State {
	case Stay3:
		e.State = Stay1
	case Idle:
		e.State = Stay1
	}
	e.mu.Unlock()
	e.signalCh <- 0
}

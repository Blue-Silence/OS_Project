package elevator

import "sync"

const (
	Upward   int = 0
	Downward     = 1
)

const (
	Idle  int = -1
	Run       = 0
	Stay1     = 1
	Stay2     = 2
	Stay3     = 3
)

type ECB struct {
	State int
	Dir   int
	floor int
	//upTarget int[]
	//downTarget int[]
	Target         []bool
	internalButton []bool
	mu             sync.Mutex
	topFloor       int
	clockCh        chan int
	signalCh       chan int
}

func MakeECB(floors int) ECB {
	var e ECB
	e.State = Idle
	e.floor = 0
	e.Dir = Upward
	e.Target = make([]bool, floors)
	e.internalButton = make([]bool, floors)
	e.topFloor = floors - 1
	e.clockCh = make(chan int)
	e.signalCh = make(chan int)
	return e
}

func (e *ECB) insertTarget(f int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Target[f-1] = true
}

func (e *ECB) insertInternalTarget(f int) {
	e.mu.Lock()
	//defer e.singalCh <- 0
	//defer e.mu.Unlock()
	e.Target[f-1] = true
	e.internalButton[f-1] = true

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

}

func (e *ECB) targetCount() (int, int) {
	upCount := 0
	downCount := 0
	for i, v := range e.Target {
		switch {
		case !v:
		case i > e.floor:
			upCount++
		case i < e.floor:
			downCount++
		}
	}
	return upCount, downCount
}

func (e *ECB) stateForwardIdle() {
	upCount, downCount := e.targetCount()

	switch {
	case e.Target[e.floor] == true:
		e.Target[e.floor] = false
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
}

func (e *ECB) stateForwardRun() {
	switch e.Dir {
	case Upward:
		e.floor++
		switch {
		case e.Target[e.floor]:
			e.State = Stay1
		case e.floor == e.topFloor:
			e.State = Stay3
		}
	case Downward:
		e.floor--
		switch {
		case e.Target[e.floor]:
			e.State = Stay1
		case e.floor == 0:
			e.State = Stay3
		}
	}
}

func (e *ECB) stateForwardStay1() {
	e.stateToStay2()
}

func (e *ECB) stateForwardStay2() {
	switch {
	case e.Target[e.floor]:
		e.stateToStay2()
	default:
		e.State = Stay3
	}
	e.Target[e.floor] = false
	e.internalButton[e.floor] = false
}

func (e *ECB) stateForwardStay3() {
	switch {
	case e.Target[e.floor]:
		e.stateToStay2()
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
	e.Target[e.floor] = false
	e.internalButton[e.floor] = false
}

func (e *ECB) stateToStay2() {
	e.State = Stay2
	e.Target[e.floor] = false
	e.internalButton[e.floor] = false
	//And also do something to clear the external button. TO BE DONE.

}

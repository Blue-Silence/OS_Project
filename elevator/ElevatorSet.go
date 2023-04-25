package elevator

import "sync"

type ESet struct {
	pannel    Pannel
	es        []ECB
	signalChs []chan int
	//elevatorNum int
	topFloor int
	mu       sync.Mutex
}

func MakeESet(floor int, elevatorNum int) ESet {
	var s ESet
	s.pannel = MakePannel(floor)
	//s.elevatorNum = elevatorNum
	s.topFloor = floor
	s.signalChs = make([]chan int, 0)
	s.es = make([]ECB, elevatorNum)
	for i, _ := range s.es {
		s.es[i] = MakeECB(floor, &s.pannel)
		s.signalChs = append(s.signalChs, s.es[i].signalCh)
	}
	return s
}

func (s *ESet) requestElevator(dir int, t int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pannel.setTarget(dir, t, true)
	//var dis [s.elevatorNum]int
	minDis := s.topFloor * 10
	chosenE := &s.es[0]
	for i, _ := range s.es {
		d := s.es[i].distanceCal(dir, t)
		if d < minDis {
			minDis = d
			chosenE = &s.es[i]
		}
	}

	chosenE.insertTarget(t)

}

package main

import (
	"log"
	"sync"
)

const (
	FIFO int = 0
	LRU  int = 1
)

const (
	PAGE_SIZE int = 10
)

type Page struct {
	vn        int
	lastUsed  int
	lastEnter int
}

type GlobalState struct {
	physicalPs  []*Page
	pageSet     []Page
	reqCounter  int
	missCounter int
	mu          sync.Mutex
	signalCh    chan ReqMsg
}

type ReqMsg struct {
	PN          int
	VN          int
	reqAddress  int
	isReplace   bool
	isHit       bool
	isReset     bool
	missCounter int
	reqCounter  int
	physicalPs  []*Page
}

func (s *GlobalState) reqAddress(addr int, replacePolicy int) bool {
	s.mu.Lock()

	vn := addr / PAGE_SIZE

	msg := ReqMsg{VN: vn, reqAddress: addr}

	s.reqCounter++
	msg.reqCounter = s.reqCounter
	msg.missCounter = s.missCounter

	for _, p := range s.physicalPs {
		if p != nil {
			p.lastEnter++
			p.lastUsed++
		}
	}

	for a, p := range s.physicalPs {
		if p != nil && p.vn == vn {
			p.lastUsed = 0
			msg.PN = a
			msg.isReplace = false
			msg.isHit = true
			msg.physicalPs = s.copyPageInfo()
			s.mu.Unlock()
			s.signalCh <- msg
			return true
		}
	} // Hit

	s.missCounter++
	msg.missCounter = s.missCounter
	for a, p := range s.physicalPs {
		if p == nil {
			s.physicalPs[a] = &s.pageSet[vn]
			s.pageSet[vn].lastUsed = 0
			s.pageSet[vn].lastEnter = 0
			msg.PN = a
			msg.isReplace = false
			msg.isHit = false
			msg.physicalPs = s.copyPageInfo()
			s.mu.Unlock()
			s.signalCh <- msg
			return false
		}
	} // Use empty PP

	pn := 0
	switch replacePolicy {
	case FIFO:
		pn = s.findPnFIFO()
	case LRU:
		pn = s.findPnLRU()
	default:
		log.Fatal("No policy matched!")
	}

	s.physicalPs[pn] = &s.pageSet[vn]
	s.pageSet[vn].lastUsed = 0
	s.pageSet[vn].lastEnter = 0
	msg.PN = pn
	msg.isReplace = true
	msg.isHit = false
	msg.physicalPs = s.copyPageInfo()
	s.mu.Unlock()
	s.signalCh <- msg
	return false

}

func (s *GlobalState) copyPageInfo() []*Page {
	var re []*Page
	//s.mu.Lock()
	//defer s.mu.Unlock()
	for _, v := range s.physicalPs {
		if v != nil {
			p := *v
			re = append(re, &p)
		} else {
			re = append(re, nil)
		}
	}
	return re
}

func (s *GlobalState) reset(physicalPN int, virtualPN int, policy int) {
	s.mu.Lock()
	s.physicalPs = make([]*Page, physicalPN)
	s.pageSet = make([]Page, virtualPN)
	for i, _ := range s.pageSet {
		s.pageSet[i].vn = i
	}
	s.reqCounter = 0
	s.missCounter = 0
	oldCh := s.signalCh
	s.signalCh = make(chan ReqMsg, 1024)
	s.mu.Unlock()

	oldCh <- ReqMsg{isReset: true}
}

func (s *GlobalState) findPnFIFO() int {
	pn := 0
	maxAge := s.physicalPs[0].lastEnter
	for a, p := range s.physicalPs {
		if p.lastEnter > maxAge {
			pn = a
			maxAge = p.lastEnter
		}
	}
	return pn
}

func (s *GlobalState) findPnLRU() int {
	pn := 0
	maxUseAge := s.physicalPs[0].lastUsed
	for a, p := range s.physicalPs {
		if p.lastUsed > maxUseAge {
			pn = a
			maxUseAge = p.lastUsed
		}
	}
	return pn
}

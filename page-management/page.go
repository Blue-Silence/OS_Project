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
	address   int
	lastUsed  int
	lastEnter int
}

type GlobalState struct {
	physicalPs    []*Page
	pageSet       []Page
	reqCounter    int
	missCounter   int
	replacePolicy int
	mu            sync.Mutex
	signalCh      chan ReqMsg
}

type ReqMsg struct {
	PN            int
	reqAdderss    int
	currentPolicy int
	isReplace     bool
	isHit         bool
}

func (s *GlobalState) reqAdderss(addr int) {
	s.mu.Lock()
	vn := addr / PAGE_SIZE

	msg := ReqMsg{reqAdderss: addr, currentPolicy: s.replacePolicy}

	for _, p := range s.physicalPs {
		if p != nil {
			p.lastEnter++
			p.lastUsed++
		}
	}

	for a, p := range s.physicalPs {
		if p != nil && p.address == vn {
			p.lastUsed = 0
			msg.PN = a
			msg.isReplace = false
			msg.isHit = true
			s.mu.Unlock()
			s.signalCh <- msg
			return
		}
	} // Hit

	for a, p := range s.physicalPs {
		if p == nil {
			s.physicalPs[a] = &s.pageSet[vn]
			s.pageSet[vn].lastUsed = 0
			s.pageSet[vn].lastEnter = 0
			msg.PN = a
			msg.isReplace = false
			msg.isHit = false
			s.mu.Unlock()
			s.signalCh <- msg
			return
		}
	} // Use empty PP

	pn := 0
	switch s.replacePolicy {
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
	s.mu.Unlock()
	s.signalCh <- msg
	return

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
			maxUseAge = p.lastEnter
		}
	}
	return pn
}

// Ispired form github.com/tj/go-spin
package main

import (
	"fmt"
	"time"
)

var (
	Box = `⠄⠆⠇⠋⠙⠸⠰⠠⠰⠸⠙⠋⠇⠆`
)

type Spinner struct {
	frames []rune
	length int
	pos    int
	done   bool
}

func (s *Spinner) Start() {
	go func() {
		for !s.done {
			fmt.Print("\r" + string(s.frames[s.pos%s.length]))
			s.pos++
			time.Sleep(70 * time.Millisecond)
		}
	}()
}

func (s *Spinner) Done() {
	fmt.Print("\r")
	s.done = true
}

func NewSpin() *Spinner {
	s := &Spinner{
		frames: []rune(Box),
	}
	s.length = len(s.frames)
	s.Start()
	return s
}

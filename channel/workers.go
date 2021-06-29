package channel

import (
	"GetBingPictures/parser"
	"os"
	"sync"
)

type Worker struct {
	In   chan int
	done func()
}

func CreateWorker(wg *sync.WaitGroup, rp, fp *os.File) Worker {
	w := Worker{
		In: make(chan int, 32),
		done: func() {
			wg.Done()
		},
	}
	go doWork(w, rp, fp)
	return w
}

func doWork(w Worker, rp, fp *os.File) {
	for i := range w.In {
		parser.Parser(i, rp, fp)
	}
	w.done()
}

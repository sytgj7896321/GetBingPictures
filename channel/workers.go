package channel

import (
	"GetBingPictures/parser"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Worker struct {
	In   chan int
	done func()
}

func CreateWorker(wg *sync.WaitGroup, fp *os.File) Worker {
	w := Worker{
		In: make(chan int, 32),
		done: func() {
			wg.Done()
		},
	}

	//go doWork(id, w, fp, logMap)
	go doWork(w, fp)
	return w
}

func doWork(w Worker, fp *os.File) {
	for i := range w.In {
		//if !logMap[i] {
		parser.Parser(i, fp)
		//}
	}
	w.done()
}

func ScannerLog(fp *os.File, overwrite bool) map[int]bool {
	var logScanner = map[int]bool{}
	scanner := bufio.NewScanner(fp)
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Read log Error %s\n", err)
		return nil
	} else {
		for scanner.Scan() {
			stringSlice := strings.Split(scanner.Text(), " ")
			id, _ := strconv.Atoi(stringSlice[3])
			if overwrite && stringSlice[4] == "Found" {
				logScanner[id] = false
			} else {
				logScanner[id] = true
			}
		}
		return logScanner
	}
}

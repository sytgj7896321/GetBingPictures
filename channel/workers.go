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

func CreateWorker(id int, wg *sync.WaitGroup, fp *os.File, logMap map[int]bool) Worker {
	w := Worker{
		In: make(chan int, 1024),
		done: func() {
			wg.Done()
		},
	}
	go DoWork(id, w, fp, logMap)
	return w
}

func DoWork(id int, w Worker, fp *os.File, logMap map[int]bool) {
	for i := range w.In {
		if !logMap[i] {
			parser.Parser(i, id, fp)
			w.done()
		}
	}
}

func ScannerLog(fp *os.File, overwrite bool) map[int]bool {
	var logScanner = map[int]bool{}
	scanner := bufio.NewScanner(fp)
	if err := scanner.Err(); err != nil {
		fmt.Fprint(os.Stderr, "Read log Error", err)
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

package main

import (
	"GetBingPictures/channel"
	"GetBingPictures/parser"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	logName     = "record.log"
	concurrency = 8
)

var (
	overwrite bool
	timeTick  int
)

func main() {
	flag.BoolVar(&overwrite, "w", false, "Overwrite Mode: skip not found, re-download found pictures)")
	flag.IntVar(&timeTick, "t", 800, "Set number of millisecond between sending http requests, require not quick than 300")
	flag.Parse()

	endNum := parser.FetchNewestId(parser.HomePage)
	err := createPath(parser.Path)
	if err != nil {
		fmt.Fprint(os.Stderr, ", Exit Programme ", err)
		return
	}
	fp, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprint(os.Stderr, ", Exit Programme ", err)
		return
	}
	defer fp.Close()
	mapLog := channel.ScannerLog(fp, overwrite)
	rateLimiter := time.Tick(time.Duration(timeTick) * time.Millisecond)
	getBingPictures(endNum, fp, mapLog, rateLimiter)

}

func getBingPictures(endNum int, fp *os.File, logMap map[int]bool, rateLimiter <-chan time.Time) {
	var wg sync.WaitGroup
	var workers [concurrency]channel.Worker
	wg.Add(endNum)
	for i := 0; i < concurrency; i++ {
		workers[i] = channel.CreateWorker(i, &wg, fp, logMap, rateLimiter)
	}

	for task := 0; task <= endNum; task++ {
		for i, worker := range workers {
			if task%(concurrency+1) == i {
				worker.In <- task
			}
		}
	}
	wg.Wait()
}

func createPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Path '" + path + "' created")
		return nil
	}
	return err
}

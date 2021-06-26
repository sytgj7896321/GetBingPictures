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
)

func main() {
	flag.BoolVar(&overwrite, "w", false, "Overwrite Mode: skip not found, re-download found pictures)")
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
	getBingPictures(endNum, fp, mapLog)
	fmt.Println("All done, Ctrl+C or 5 seconds later will exit programme")
	timeout := time.After(5 * time.Second)
	for {
		select {
		case <-timeout:
			fmt.Println("Exited")
			return
		}
	}
}

func getBingPictures(endNum int, fp *os.File, logMap map[int]bool) {
	var wg sync.WaitGroup
	var workers [concurrency]channel.Worker
	wg.Add(endNum + 1)
	for i := 0; i < concurrency; i++ {
		workers[i] = channel.CreateWorker(i, &wg, fp, logMap)
	}

	for task := 0; task <= endNum; task++ {
		for i, worker := range workers {
			if task%(concurrency+1) == i {
				worker.In <- task
			}
		}
	}
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

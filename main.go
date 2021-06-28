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
	logName    = "record.log"
	goroutines = 8
)

var (
	goroutine int
)

func main() {
	flag.IntVar(&goroutine, "c", 16, "Set how many coroutines you want to use")
	flag.Parse()

	lastNum, _ := parser.FetchLatestPageNum()
	err := createPath(parser.Path)
	if err != nil {
		fmt.Fprint(os.Stderr, err, ", exit programme\n")
		return
	}
	fp, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprint(os.Stderr, err, ", exit programme\n")
		return
	}
	defer fp.Close()

	//mapLog := channel.ScannerLog(fp, overwrite)
	//Engine(lastNum, fp, mapLog)
	Engine(lastNum, fp)
	timeout := time.After(3 * time.Second)
	for {
		select {
		case <-timeout:
			fmt.Println("Exited")
			return
		}
	}
}

func Engine(lastNum int, fp *os.File) {
	var wg sync.WaitGroup
	var workers = make([]channel.Worker, goroutines)
	fmt.Println(lastNum)
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		//workers[i] = channel.CreateWorker(i, &wg, fp, logMap)
		workers[i] = channel.CreateWorker(i, &wg, fp)

	}
	for task := 0; task <= lastNum; task++ {
		for i, worker := range workers {
			if task%(goroutines+1) == i {
				worker.In <- task
			}
		}
	}
	for i := range workers {
		close(workers[i].In)
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

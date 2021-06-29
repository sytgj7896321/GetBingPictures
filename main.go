package main

import (
	"GetBingPictures/channel"
	"GetBingPictures/fetcher"
	"GetBingPictures/parser"
	"flag"
	"fmt"
	"os"
	"sync"
)

const (
	recordLog = "DownloadRecord.log"
	errLog    = "Error.log"
)

var goroutines int

func main() {
	flag.IntVar(&goroutines, "c", 4, "Set how many coroutines you want to use like -c 8")
	flag.StringVar(&fetcher.ProxyAdd, "p", "", "Set http proxy address like -p \"http://127.0.0.1:10809\"")
	flag.Parse()

	lastNum, _ := parser.FetchLatestPageNum()
	err := createPath(parser.Path)
	if err != nil {
		fmt.Fprint(os.Stderr, err, ", exit programme\n")
		return
	}
	rp, err := os.OpenFile(recordLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprint(os.Stderr, err, ", exit programme\n")
		return
	}
	defer rp.Close()
	fp, err := os.OpenFile(errLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprint(os.Stderr, err, ", exit programme\n")
		return
	}
	defer fp.Close()

	parser.ScannerRecord(rp)
	//Engine(lastNum, fp, mapLog)
	Engine(lastNum, rp, fp)
	//timeout := time.After(3 * time.Second)
	//for {
	//	select {
	//	case <-timeout:
	//		fmt.Println("Exited")
	//		return
	//	}
	//}
}

func Engine(lastNum int, rp, fp *os.File) {
	var wg sync.WaitGroup
	var workers = make([]channel.Worker, goroutines)
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		workers[i] = channel.CreateWorker(&wg, rp, fp)

	}
	for task := 1; task <= lastNum; task++ {
		for i, worker := range workers {
			if task%goroutines == i {
				worker.In <- task
			}
		}
	}
	//for  _, worker := range workers {
	//	close(worker.In)
	//}
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
		fmt.Println("Directory '" + path + "' created")
		return nil
	}
	return err
}

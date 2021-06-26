package main

import (
	"GetBingPictures/parser"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
)

type worker struct {
	in   chan int
	done func()
}

const (
	homePage       = "https://wallpaperhub.app/"
	target         = "https://wallpaperhub.app/wallpapers/"
	path           = "wallpapers"
	workerTotalNum = 20
)

var (
	start   = 9244
	end     = regexp.MustCompile(`<a href="/wallpapers/([0-9]+)">View</a>`)
	picName = regexp.MustCompile(`<title data-react-helmet="true">(.+) \| Wallpapers \| WallpaperHub</title>`)
	picUrl  = regexp.MustCompile(`<img src="(https://cdn.wallpaperhub.app/cloudcache/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]{40}\.jpg)"/>`)
)

func main() {
	endNum := parser.FetchNewestId(homePage, end)
	err := createPath(path)
	if err != nil {
		fmt.Fprint(os.Stderr, ", Exit Programme", err)
		return
	}
	fp, err := os.OpenFile("record.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprint(os.Stderr, ", Exit Programme", err)
		return
	}

	getBingPictures(endNum, fp)

	defer fp.Close()
	fmt.Println("Press any key to exit")
	select {}
}

func getBingPictures(endNum int, fp *os.File) {
	var wg sync.WaitGroup
	var workers [workerTotalNum]worker
	wg.Add(endNum)
	for i := 0; i < workerTotalNum; i++ {
		workers[i] = createWorker(i, &wg, fp)
	}

	for task := start; task <= endNum; task++ {
		for i, worker := range workers {
			if task%(workerTotalNum+1) == i {
				worker.in <- task
			}
		}
	}
	wg.Wait()
}

func createWorker(id int, wg *sync.WaitGroup, fp *os.File) worker {
	w := worker{
		in: make(chan int, 1024),
		done: func() {
			wg.Done()
		},
	}
	go doWork(id, w, fp)
	return w
}

func doWork(id int, w worker, fp *os.File) {
	for i := range w.in {
		parser.Parser(i, id, target+strconv.Itoa(i), path, picName, picUrl, fp)
		w.done()
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

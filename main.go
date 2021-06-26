package main

import (
	"GetBingPictures/parser"
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	logName        = "record.log"
	workerTotalNum = 20
)

var (
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
	fp, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprint(os.Stderr, ", Exit Programme", err)
		return
	}
	defer fp.Close()
	mapLog := scannerLog(fp)
	getBingPictures(endNum, fp, mapLog)

}

func getBingPictures(endNum int, fp *os.File, logMap map[int]bool) {
	var wg sync.WaitGroup
	var workers [workerTotalNum]worker
	wg.Add(endNum)
	for i := 0; i < workerTotalNum; i++ {
		workers[i] = createWorker(i, &wg, fp, logMap)
	}

	for task := 0; task <= endNum; task++ {
		for i, worker := range workers {
			if task%(workerTotalNum+1) == i {
				worker.in <- task
			}
		}
	}
	wg.Wait()
}

func createWorker(id int, wg *sync.WaitGroup, fp *os.File, logMap map[int]bool) worker {
	w := worker{
		in: make(chan int, 1024),
		done: func() {
			wg.Done()
		},
	}
	go doWork(id, w, fp, logMap)
	return w
}

func doWork(id int, w worker, fp *os.File, logMap map[int]bool) {
	for i := range w.in {
		if !logMap[i] {
			parser.Parser(i, id, target+strconv.Itoa(i), path, picName, picUrl, fp)
			w.done()
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

func scannerLog(fp *os.File) map[int]bool {
	var logScanner = map[int]bool{}
	scanner := bufio.NewScanner(fp)
	if err := scanner.Err(); err != nil {
		fmt.Fprint(os.Stderr, "Read log Error", err)
		return nil
	} else {
		for scanner.Scan() {
			total := strings.Split(scanner.Text(), " ")
			id, _ := strconv.Atoi(total[3])
			logScanner[id] = true
		}
		return logScanner
	}
}

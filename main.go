package main

import (
	"GetBingPictures/channel"
	"GetBingPictures/lib"
	"GetBingPictures/parser"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
)

const (
	homePage = "https://wallpaperhub.app/"
	target   = "https://wallpaperhub.app/wallpapers/"
	path     = "wallpapers"
)

var (
	wg      sync.WaitGroup
	start   = 0
	end     = regexp.MustCompile(`<a href="/wallpapers/([0-9]+)">View</a>`)
	picName = regexp.MustCompile(`<title data-react-helmet="true">(.+) \| Wallpapers \| WallpaperHub</title>`)
	picUrl  = regexp.MustCompile(`<img src="(https://cdn.wallpaperhub.app/cloudcache/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]{40}\.jpg)"/>`)
)

func main() {
	parser.FetchNewestId(homePage, end)
	err := mylib.CreatePath(path)
	if err != nil {
		fmt.Fprint(os.Stderr, ", Exit Programme", err)
		return
	}

	getBingPictures()

	fmt.Scanf("%s", "Press any key to exit")
}

func getBingPictures() {
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		channels[i] = channel.CreateWorker(i, worker)
	}
}

func worker(id int, ch chan int) {
	for i := range ch {
		parser.Parser(i, target+strconv.Itoa(id), path, picName, picUrl)
	}
	wg.Done()
}

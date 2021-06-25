package main

import (
	"GetBingPictures/fetcher"
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
	wg          sync.WaitGroup
	start       = 0
	end         = regexp.MustCompile(`<a href="/wallpapers/([0-9]+)">View</a>`)
	picName     = regexp.MustCompile(`<title data-react-helmet="true">(.+) \| Wallpapers \| WallpaperHub</title>`)
	picUrl      = regexp.MustCompile(`<img src="(https://cdn.wallpaperhub.app/cloudcache/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]{40}\.jpg)"/>`)
)

func main() {
	result, _ := fetcher.Fetch(homePage)
	subMatch := end.FindSubmatch(result)
	endNum, _ := strconv.Atoi(string(subMatch[1]))
	fmt.Printf("Newest Wallpaper ID is %d\n", endNum)

	err := mylib.CreatePath(path)
	if err != nil {
		fmt.Fprint(os.Stderr, ", Exit Programme", err)
		return
	}


	fmt.Scanf("%s", "Blocking, press any key to exit")
}

func createWorker(id int) chan<- int {
	for i := range ch {
		parser.Parser(i, target+strconv.Itoa(id), path, picName, picUrl)
	}
	wg.Done()
	return nil
}







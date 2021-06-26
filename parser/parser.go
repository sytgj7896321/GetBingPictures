package parser

import (
	"GetBingPictures/fetcher"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const (
	HomePage = "https://wallpaperhub.app/"
	Path     = "wallpapers"
)

var (
	end     = regexp.MustCompile(`<a href="/wallpapers/([0-9]+)">View</a>`)
	picName = regexp.MustCompile(`<title data-react-helmet="true">(.+) \| Wallpapers \| WallpaperHub</title>`)
	picUrl  = regexp.MustCompile(`<img src="(https://cdn.wallpaperhub.app/cloudcache/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]/[0-9a-z]{40}\.jpg)"/>`)
)

func Parser(pid, wid int, fp *os.File) {
	log.SetOutput(fp)
	log.SetPrefix("[GetBingTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Printf("Worker %d received Task %d, and begin fetching\n", wid, pid)
	result, _ := fetcher.Fetch(HomePage + Path + "/" + strconv.Itoa(pid))
	subMatch1 := picName.FindSubmatch(result)
	subMatch2 := picUrl.FindSubmatch(result)
	if subMatch1 == nil || subMatch2 == nil {
		fmt.Printf("No wallpaper with ID %d found\n", pid)
		log.Printf("%d Not\u00a0found\n", pid)
		return
	}
	subMatch1[1] = append(subMatch1[1], ".jpg"...)

	fmt.Printf("Find and begin download: %s\n", subMatch1[1])
	resp, err := http.Get(string(subMatch2[1]))
	if err != nil {
		fmt.Fprint(os.Stderr, "Get Image Error", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprint(os.Stderr, "IO Read Error", err)
		return
	}
	defer resp.Body.Close()

	err = ioutil.WriteFile(Path+"/"+string(subMatch1[1]), data, 0755)
	if err != nil {
		fmt.Fprint(os.Stderr, "IO Write Error", err)
		return
	}
	fmt.Printf("ID %d %s download completed\n", pid, string(subMatch1[1]))
	log.Printf("%d Found %s\n", pid, string(subMatch1[1]))

}

func FetchNewestId(homePage string) int {
	result, _ := fetcher.Fetch(homePage)
	subMatch := end.FindSubmatch(result)
	endNum, _ := strconv.Atoi(string(subMatch[1]))
	fmt.Printf("Newest Wallpaper ID is %d\n", endNum)
	return endNum
}

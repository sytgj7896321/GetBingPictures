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

func Parser(pid, wid int, url, path string, picName, picUrl *regexp.Regexp, fp *os.File) {
	log.SetOutput(fp)
	log.SetPrefix("[GetBingTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Printf("Worker %d received Task %d, and begin fetching\n", wid, pid)
	result, _ := fetcher.Fetch(url)
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

	err = ioutil.WriteFile(path+"/"+string(subMatch1[1]), data, 0755)
	if err != nil {
		fmt.Fprint(os.Stderr, "IO Write Error", err)
		return
	}
	fmt.Printf("ID %d %s download completed\n", pid, string(subMatch1[1]))
	log.Printf("%d Downloaded %s\n", pid, string(subMatch1[1]))

}

func FetchNewestId(homePage string, end *regexp.Regexp) int {
	result, _ := fetcher.Fetch(homePage)
	subMatch := end.FindSubmatch(result)
	endNum, _ := strconv.Atoi(string(subMatch[1]))
	fmt.Printf("Newest Wallpaper ID is %d\n", endNum)
	return endNum
}

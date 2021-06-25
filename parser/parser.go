package parser

import (
	"GetBingPictures/fetcher"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var rateLimiter = time.Tick(100 * time.Millisecond)

func Parser(id int, url, path string, picName, picUrl *regexp.Regexp) {
	<-rateLimiter
	fmt.Println("Fetching:", url)
	result, _ := fetcher.Fetch(url)
	subMatch1 := picName.FindSubmatch(result)
	subMatch2 := picUrl.FindSubmatch(result)
	if subMatch1 == nil || subMatch2 == nil {
		fmt.Printf("No wallpaper with ID %d found\n", id)
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
	fmt.Printf("%s download completed\n", string(subMatch1[1]))

}

func FetchNewestId(homePage string, end *regexp.Regexp) {
	result, _ := fetcher.Fetch(homePage)
	subMatch := end.FindSubmatch(result)
	endNum, _ := strconv.Atoi(string(subMatch[1]))
	fmt.Printf("Newest Wallpaper ID is %d\n", endNum)
}
